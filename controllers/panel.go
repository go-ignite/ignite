package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ignite/models"
)

func (router *MainRouter) PanelIndexHandler(c *gin.Context) {
	userID, exists := c.Get("userId")

	if exists {
		user := new(models.User)
		router.db.Id(userID).Get(user)

		c.HTML(http.StatusOK, "panel.html", gin.H{
			"username": user.Username,
		})
	}

	c.HTML(http.StatusOK, "panel.html", nil)
}

func (router *MainRouter) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userId")
	session.Save()

	c.Redirect(http.StatusFound, "/")
}
