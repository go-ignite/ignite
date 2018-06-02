package handler

import (
	"context"
	"net/http"

	pb "github.com/go-ignite/ignite-agent/protos"
	"github.com/go-ignite/ignite/agent"
	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (ah *AdminHandler) AddNodeHandler(c *gin.Context) {
	nodeEntity := struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}{}
	if err := c.BindJSON(&nodeEntity); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResp(err.Error()))
		return
	}
	logrus.WithField("nodeEntity", nodeEntity).Debug()

	agentClient, err := agent.Dial(nodeEntity.Address)
	if err != nil {
		c.JSON(http.StatusOK, models.NewErrorResp("不能与节点建立连接！"))
		return
	}

	req := &pb.GeneralRequest{Token: c.GetString("token")}
	if _, err = agentClient.Init(context.Background(), req); err != nil {
		logrus.WithField("err", err).Error("agent init error")
		c.JSON(http.StatusOK, models.NewErrorResp("节点初始化失败！"))
		return
	}

	node := &db.Node{
		Name:    nodeEntity.Name,
		Address: nodeEntity.Address,
	}
	affected, err := db.CreateNode(node)
	if err != nil || affected == 0 {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"affected": affected,
		}).Error("create node in db error")
		c.JSON(http.StatusOK, models.NewErrorResp("添加节点失败！"))
		return
	}
	c.JSON(http.StatusOK, models.NewSuccessResp(node))
}
