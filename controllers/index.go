package controllers

import (
	"fmt"
	"ignite/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"

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
	router.db.Where("invite_code = ? && available = 1", inviteCode).Get(iv)

	if iv.Id == 0 {
		fmt.Println("Invalid invite code!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Invalid invite code!"})
		return
	}

	user := new(models.User)
	count, _ := router.db.Where("username = ?", username).Count(user)

	if count > 0 {
		fmt.Println("Username duplicated!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Username is duplicated!"})
		return
	}

	//Create user account
	session := router.db.NewSession()
	defer session.Close()

	session.Begin()

	//1.Create user account
	user.Username = username
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	user.HashedPwd = hashedPass
	user.InviteCode = iv.InviteCode
	user.PackageLimit = iv.PackageLimit

	affected, _ := session.Insert(user)

	if affected == 0 {
		session.Rollback()
		fmt.Println("Failed to create user account!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Failed to create user account!"})
		return
	}

	//2.Set invite code as used status
	iv.Available = false
	affected, _ = session.Id(iv.Id).Cols("available").Update(iv)

	if affected == 0 {
		session.Rollback()
		fmt.Println("Failed to create user account!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Failed to create user account!"})
		return
	}

	if err := session.Commit(); err != nil {
		fmt.Println("Failed to create user account!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Failed to create user account!"})
		return
	}

	fmt.Printf("User %s with invite code: %s", username, inviteCode)
	c.JSON(http.StatusOK, &models.Response{Success: true, Message: "Success!"})
}
