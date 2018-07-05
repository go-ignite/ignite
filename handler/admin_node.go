package handler

import (
	"context"
	"net/http"
	"strconv"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/state"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/db/api"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

func (ah *AdminHandler) AddNode(c *gin.Context) {
	nodeEntity := &models.AddNodeReq{}
	if err := c.BindJSON(nodeEntity); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp(err.Error()))
		return
	}
	ah.logger.WithFields(logrus.Fields{
		"name":       nodeEntity.Name,
		"comment":    nodeEntity.Comment,
		"address":    nodeEntity.Address,
		"connection": nodeEntity.Connection,
		"port_from":  nodeEntity.PortFrom,
		"port_to":    nodeEntity.PortTo,
	}).Debug()

	agentClient, err := agent.Dial(nodeEntity.Address)
	if err != nil {
		ah.logger.WithError(err).Error("agent dial error")
		c.JSON(http.StatusOK, models.NewErrorResp("不能与节点建立连接！"))
		return
	}

	req := &pb.GeneralRequest{
		Token: c.GetString("token"),
	}
	if _, err = agentClient.Init(context.Background(), req); err != nil {
		agentClient.Close()
		ah.logger.WithError(err).Error("agent init error")
		c.JSON(http.StatusOK, models.NewErrorResp("节点初始化失败！"))
		return
	}

	node := &db.Node{}
	copier.Copy(node, nodeEntity)
	affected, err := api.NewAPI().UpsertNode(node)
	if err != nil || affected == 0 {
		ah.logger.WithFields(logrus.Fields{
			"error":    err,
			"affected": affected,
		}).Error("add node error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("添加节点失败！"))
		return
	}

	go state.GetLoader().AddNode(node.Id, state.NewNodeStatus(node, agentClient, true, nil))
	nodeResp := &models.NodeResp{}
	copier.Copy(nodeResp, node)
	nodeResp.Available = true
	c.JSON(http.StatusOK, models.NewSuccessResp(nodeResp))
}

func (ah *AdminHandler) ListNodes(c *gin.Context) {
	nodes, err := api.NewAPI().GetAllNodes()
	if err != nil {
		ah.logger.WithError(err).Error("list nodes error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取节点列表失败！"))
		return
	}
	var nodeResps []*models.NodeResp
	for _, node := range nodes {
		nodeResp := &models.NodeResp{}
		copier.Copy(nodeResp, node)
		nodeResp.Available = state.GetLoader().GetNodeAvailable(node.Id)
		nodeResps = append(nodeResps, nodeResp)
	}
	c.JSON(http.StatusOK, models.NewSuccessResp(nodeResps))
}

func (ah *AdminHandler) DeleteNode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp("id is invalid"))
		return
	}
	affected, err := api.NewAPI().DeleteNode(id)
	if err != nil || affected == 0 {
		ah.logger.WithFields(logrus.Fields{
			"error":    err,
			"affected": affected,
		}).Error("delete node error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("移除节点失败！"))
		return
	}
	go state.GetLoader().DelNode(id)
	c.JSON(http.StatusOK, models.NewSuccessResp(nil))
}

func (ah *AdminHandler) UpdateNode(c *gin.Context) {
	nodeEntity := &models.UpdateNodeReq{}
	if err := c.BindJSON(nodeEntity); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp(err.Error()))
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp("id is invalid"))
		return
	}
	ah.logger.WithFields(logrus.Fields{
		"id":         id,
		"name":       nodeEntity.Name,
		"comment":    nodeEntity.Comment,
		"connection": nodeEntity.Connection,
		"port_from":  nodeEntity.PortFrom,
		"port_to":    nodeEntity.PortTo,
	}).Debug("update node")

	node := &db.Node{Id: id}
	copier.Copy(node, nodeEntity)
	if _, err = api.NewAPI().UpsertNode(node); err != nil {
		ah.logger.WithError(err).Error("update node error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("更新节点失败！"))
		return
	}
	go func() {
		ns := state.GetLoader().GetNode(id)
		ns.Lock()
		defer ns.Unlock()
		copier.Copy(ns.Node, nodeEntity)
	}()
	c.JSON(http.StatusOK, models.NewSuccessResp(nil))
}
