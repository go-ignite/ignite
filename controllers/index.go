package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/models"
)

func (router *MainRouter) IndexHandler(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("userId")
	var uInfo *models.UserInfo
	if v != nil {
		if uId, ok := v.(int64); ok {
			uInfo = &models.UserInfo{
				Id: uId,
			}
		}
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"uInfo": uInfo,
	})
}

func (router *MainRouter) LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("password")

	user := new(models.User)
	router.db.Where("username = ?", username).Get(user)

	if user.Id == 0 {
		fmt.Println("User doesn't exist!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "User doesn't exist!"})
		return
	}

	if bcrypt.CompareHashAndPassword(user.HashedPwd, []byte(pwd)) != nil {
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Username or password is incorrect!"})
		return
	}

	fmt.Println("Come here...")
	fmt.Println("userId is:", user.Id)

	session := sessions.Default(c)
	session.Set("userId", user.Id)
	session.Save()

	c.JSON(http.StatusOK, &models.Response{Success: true, Message: "Success!"})
}

func (router *MainRouter) SignupHandler(c *gin.Context) {
	inviteCode := c.PostForm("invite-code")
	username := c.PostForm("username")
	pwd := c.PostForm("password")
	confirmPwd := c.PostForm("confirm-password")

	matched, _ := regexp.MatchString("^[a-zA-Z0-9][a-zA-Z0-9_.-]+$", username)

	if !matched {
		fmt.Println("Username is invalid!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Username is invalid!"})
		return
	}

	if pwd != confirmPwd {
		fmt.Println("passwords not match!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Passwords don't match!"})
		return
	}

	iv := new(models.InviteCode)
	router.db.Where("invite_code = ? AND available = 1", inviteCode).Get(iv)

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
	trans := router.db.NewSession()
	defer trans.Close()

	trans.Begin()

	//1.Create user account
	user.Username = username
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	user.HashedPwd = hashedPass
	user.InviteCode = iv.InviteCode
	user.PackageLimit = iv.PackageLimit
	user.Expired = time.Now().AddDate(0, iv.AvailableLimit, 0)

	affected, _ := trans.Insert(user)

	if affected == 0 {
		trans.Rollback()
		fmt.Println("Failed to create user account!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Failed to create user account!"})
		return
	}

	//2.Set invite code as used status
	iv.Available = false
	affected, _ = trans.Id(iv.Id).Cols("available").Update(iv)

	if affected == 0 {
		trans.Rollback()
		fmt.Println("Failed to create user account!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Failed to create user account!"})
		return
	}

	if err := trans.Commit(); err != nil {
		fmt.Println("Failed to create user account!")
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Failed to create user account!"})
		return
	}

	session := sessions.Default(c)
	session.Set("userId", user.Id)
	session.Save()

	fmt.Printf("User %s with invite code: %s", username, inviteCode)
	c.JSON(http.StatusOK, &models.Response{Success: true, Message: "Success!"})
}
