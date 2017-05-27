package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (router *MainRouter) IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tpl", gin.H{
		"title": "Main website",
	})
}

func (router *MainRouter) LoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tpl", gin.H{
		"title": "Main website",
	})
}

func (router *MainRouter) SignupHandler(c *gin.Context) {
	inviteCode := c.PostForm("invite-code")
	username := c.PostForm("username")
	pwd := c.PostForm("password")
	confirmPwd := c.PostForm("confirm-password")

	if pwd != confirmPwd {
		fmt.Println("password not match!")
	}

	fmt.Printf("User %s with invite code: %s", username, inviteCode)

	c.HTML(http.StatusOK, "index.tpl", gin.H{
		"title": "Main website",
	})
}
