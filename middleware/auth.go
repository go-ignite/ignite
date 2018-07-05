package middleware

import (
	"fmt"
	"time"

	"github.com/go-ignite/ignite/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func Auth(admin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(config.C.App.Secret))
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
		if (admin && id != -1) || (!admin && id <= 0) {
			c.AbortWithError(401, fmt.Errorf("token auth error"))
			return
		}

		c.Set("id", claims["id"])
		c.Set("token", token.Raw)
	}
}
