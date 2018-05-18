package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/config"
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

	token, err := createToken(config.C.Auth.Secret, user.Id)
	resp := models.Response{}
	if err != nil {
		resp.Success = false
		resp.Message = "Could not generate token"
		c.JSON(http.StatusInternalServerError, &resp)
		return
	}

	resp.Success = true
	resp.Message = "success"
	resp.Data = token
	c.JSON(http.StatusOK, &resp)
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

	token, err := createToken(config.C.Auth.Secret, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, &models.Response{Success: false, Message: "Could not generate token"})
		return
	}

	fmt.Printf("User %s with invite code: %s", username, inviteCode)
	c.JSON(http.StatusOK, &models.Response{Success: true, Message: "Success!", Data: token})
}

func createToken(secret string, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenStr, nil
}
