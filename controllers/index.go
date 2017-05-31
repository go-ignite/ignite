package controllers

import (
	"fmt"
	"ignite/models"
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
		fmt.Println("passwords not match!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Passwords don't match!"})
	}

	iv := new(models.InviteCode)
	count, _ := router.db.Where("invite_code = ?", inviteCode).Count(iv)

	if count == 0 {
		fmt.Println("Invalid invite code!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Invalid invite code!"})
		return
	}

	user := new(models.User)
	count, _ = router.db.Where("username = ?", username).Count(user)

	if count > 0 {
		fmt.Println("Username duplicated!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Username is duplicated!"})
		return
	}

	fmt.Printf("User %s with invite code: %s", username, inviteCode)
	c.JSON(http.StatusOK, &models.Response{Success: true, Message: "Success!"})
}
