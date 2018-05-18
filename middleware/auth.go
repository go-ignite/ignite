package middleware

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(secret))
			return b, nil
		})
		if err != nil {
			c.AbortWithError(401, err)
		}
		if !token.Valid {
			c.AbortWithError(401, fmt.Errorf("token valid"))
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("id", claims["id"])
	}
}
