package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/api"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/state"
)

// 1 -> cannot connect to node
// 2 -> cannot init node
// 3 -> node name already exists
func (ah *AdminHandler) AddNode(c *gin.Context) {
	req := new(api.AddNodeRequest)
	if err := c.ShouldBind(req); err != nil {
		ah.ErrJSON(c, http.StatusBadRequest, err)
		return
	}
	ah.logger.WithField("req", fmt.Sprintf("%+v", req)).Debug()
	if models.CheckIfNodeNameExist(req.Name) {
		ah.ErrJSON(c, http.StatusBadRequest, models.ErrDuplicateNodeName, 3)
		return
	}

	agentClient, err := agent.Dial(req.RequestAddress)
	if err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "connect node error"), 1)
		return
	}

	gr := &protos.GeneralRequest{Token: ah.getToken(c)}
	if _, err = agentClient.Init(context.Background(), gr); err != nil {
		agentClient.Close()
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "init node error"), 2)
		return
	}

	node := models.NewNode(req.Name, req.Comment, req.RequestAddress, req.ConnectionAddress, req.PortFrom, req.PortTo)
	if err := node.Create(); err != nil {
		if err == models.ErrDuplicateNodeName {
			ah.ErrJSON(c, http.StatusBadRequest, models.ErrDuplicateNodeName, 3)
		} else {
			ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "create node error"))
		}
		return
	}

	go state.GetLoader().AddNodeStatus(state.NewNodeStatus(node, agentClient, nil))

	resp := new(api.Node)
	if err := copier.Copy(resp, node); err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (ah *AdminHandler) GetAllNodes(c *gin.Context) {
	nodes, err := models.GetAllNodes()
	if err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get nodes error"))
		return
	}

	resp := make([]*api.Node, 0)
	for _, node := range nodes {
		n := new(api.Node)
		if err := ah.copy(n, node); err != nil {
			ah.ErrJSON(c, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, n)
	}

	c.JSON(http.StatusOK, resp)
}

func (ah *AdminHandler) DeleteNode(c *gin.Context) {
	id := c.Param("id")
	ah.logger.WithField("id", id).Debug("DeleteNode")

	if err := models.DeleteNodeByID(id); err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "delete node error"))
		return
	}

	go state.GetLoader().DeleteNodeStatus(id)
	c.JSON(http.StatusNoContent, nil)
}

// 1 -> there are services outside the port range
func (ah *AdminHandler) UpdateNode(c *gin.Context) {
	req := new(api.UpdateNodeRequest)
	if err := c.ShouldBind(req); err != nil {
		ah.ErrJSON(c, http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")
	ah.logger.WithFields(logrus.Fields{
		"id":  id,
		"req": fmt.Sprintf("%+v", req),
	}).Debug("UpdateNode")

	count, err := models.GetServiceCountByNodeIDAndPortRange(id, req.PortFrom, req.PortTo)
	if err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "get service count error"))
		return
	}
	if count > 0 {
		ah.ErrJSON(c, http.StatusBadRequest, fmt.Errorf("there are %d services outside the port range", count), 1)
		return
	}

	node := &models.Node{ID: id}
	if err := ah.copy(node, req); err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, err)
		return
	}

	if err = node.Save(); err != nil {
		ah.ErrJSON(c, http.StatusInternalServerError, errors.Wrap(err, "save node error"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
