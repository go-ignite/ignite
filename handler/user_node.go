package handler

import (
	"net/http"

	"github.com/go-ignite/ignite/db"
	"github.com/go-ignite/ignite/models"

	"github.com/gin-gonic/gin"
)

func (uh *UserHandler) ListNodes(c *gin.Context) {
	nodes, err := db.GetAllNodes()
	if err != nil {
		uh.WithError(err).Error("get node list error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("获取节点列表失败！"))
		return
	}
	c.JSON(http.StatusOK, models.NewSuccessResp(nodes))
}
