package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/models"
	"github.com/go-ignite/ignite/utils"
	"github.com/sirupsen/logrus"
)

type AdminHandler struct {
	logger *logrus.Logger
}

func NewAdminHandler(l *logrus.Logger) *AdminHandler {
	return &AdminHandler{
		logger: l,
	}
}

func (ah *AdminHandler) PanelLoginHandler(c *gin.Context) {
	loginEntity := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BindJSON(&loginEntity); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("parse json data error"))
		return
	}
	ah.logger.WithField("loginEntity", loginEntity).Debug()

	if loginEntity.Username != config.C.Admin.Username || loginEntity.Password != config.C.Admin.Password {
		c.JSON(http.StatusOK, models.NewErrorResp("用户名或密码错误！"))
		return
	}
	// Create the token
	token, err := utils.CreateToken(config.C.App.Secret, -1)
	if err != nil {
		logrus.WithField("err", err).Error("generate token error")
		c.JSON(http.StatusInternalServerError, models.NewErrorResp("登录失败！"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResp(token))
	return
}
