package middleware

import (
	"fmt"
	"time"

	"github.com/go-ignite/ignite/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	isAdmin bool
}

func NewUserAuthHandler() *AuthHandler {
	return &AuthHandler{
		isAdmin: false,
	}
}

func NewAdminAuthHandler() *AuthHandler {
	return &AuthHandler{
		isAdmin: true,
	}
}

func (ah *AuthHandler) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			b := []byte(config.C.App.Secret)
			return b, nil
		})
		if err != nil {
			c.AbortWithError(401, err)
			return
		}
		if !token.Valid {
			c.AbortWithError(401, fmt.Errorf("token is invalid"))
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			c.AbortWithError(401, fmt.Errorf("token is expired"))
			return
		}
		id, ok := claims["id"].(float64)
		if !ok {
			c.AbortWithError(401, fmt.Errorf("token'id is invalid"))
			return
		}
		if (ah.isAdmin && id != -1) || (!ah.isAdmin && id <= 0) {
			c.AbortWithError(401, fmt.Errorf("token auth error"))
			return
		}

		c.Set("id", claims["id"])
		c.Set("token", token.Raw)
	}
}
