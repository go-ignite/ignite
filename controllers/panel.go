package controllers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (router *MainRouter) PanelIndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "panel.tpl", nil)
}

func (router *MainRouter) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userId")
	session.Save()

	c.HTML(http.StatusOK, "index.tpl", gin.H{
		"title": "Main website",
	})
}
