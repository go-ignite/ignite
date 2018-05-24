package controllers

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/models"
)

func (router *MainRouter) PanelLoginHandler(c *gin.Context) {
	loginEntity := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BindJSON(&loginEntity); err != nil {
		resp := models.Response{Success: false, Message: "Could not parse username & password..."}
		c.JSON(http.StatusInternalServerError, &resp)
		return
	}

	fmt.Println("username:", loginEntity.Username)
	fmt.Println("pwd:", loginEntity.Password)

	if loginEntity.Username == config.C.Admin.Username && loginEntity.Password == config.C.Admin.Password {
		// Create the token
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		// Set some claims
		token.Claims = jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		}
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte(config.C.Secret.Admin))
		resp := models.Response{}

		if err != nil {
			resp.Success = false
			resp.Message = "Could not generate token"
			c.JSON(http.StatusInternalServerError, &resp)
			return
		}

		resp.Success = true
		resp.Message = "success"
		resp.Data = tokenString
		c.JSON(http.StatusOK, &resp)
		return
	}

	resp := models.Response{Success: false, Message: "Username of password is wrong!"}
	c.JSON(http.StatusOK, &resp)
}
