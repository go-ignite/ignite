package handler

import (
	"context"
	"net/http"
	"strconv"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (ah *AdminHandler) AddNode(c *gin.Context) {
	nodeEntity := struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}{}
	if err := c.BindJSON(&nodeEntity); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResp(err.Error()))
		return
	}
	ah.WithFields(logrus.Fields{
		"name":    nodeEntity.Name,
		"address": nodeEntity.Address,
	}).Debug()

	agentClient, err := agent.Dial(nodeEntity.Address)
	if err != nil {
		ah.WithError(err).Error("agent dial error")
		c.JSON(http.StatusOK, models.NewErrorResp("不能与节点建立连接！"))
		return
	}

	req := &pb.GeneralRequest{
		Token: c.GetString("token"),
	}
	if _, err = agentClient.Init(context.Background(), req); err != nil {
		ah.WithError(err).Error("agent init error")
		c.JSON(http.StatusOK, models.NewErrorResp("节点初始化失败！"))
		return
	}

	node := &db.Node{
		Name:    nodeEntity.Name,
		Address: nodeEntity.Address,
	}
	affected, err := db.UpsertNode(node)
	if err != nil || affected == 0 {
		ah.WithFields(logrus.Fields{
			"error":    err,
			"affected": affected,
		}).Error("add node error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("添加节点失败！"))
		return
	}
	c.JSON(http.StatusOK, models.NewSuccessResp(node))
}

func (ah *AdminHandler) ListNodes(c *gin.Context) {
	nodes, err := db.GetAllNodes()
	if err != nil {
		ah.WithError(err).Error("list nodes error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取节点列表失败！"))
		return
	}
	c.JSON(http.StatusOK, models.NewSuccessResp(nodes))
}

func (ah *AdminHandler) DeleteNode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResp("id is invalid"))
		return
	}
	affected, err := db.DeleteNode(id)
	if err != nil || affected == 0 {
		ah.WithFields(logrus.Fields{
			"error":    err,
			"affected": affected,
		}).Error("delete node error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("移除节点失败！"))
		return
	}
	c.JSON(http.StatusOK, models.NewSuccessResp(nil))
}

func (ah *AdminHandler) UpdateNode(c *gin.Context) {
	nodeEntity := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&nodeEntity); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResp(err.Error()))
		return
	}
	ah.WithField("name", nodeEntity.Name).Debug()

	node := &db.Node{
		Name: nodeEntity.Name,
	}
	affected, err := db.UpsertNode(node)
	if err != nil || affected == 0 {
		ah.WithFields(logrus.Fields{
			"error":    err,
			"affected": affected,
		}).Error("update node error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("更新节点失败！"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResp(nil))
}
