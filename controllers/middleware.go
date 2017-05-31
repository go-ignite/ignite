package controllers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ValidateSession() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		v := session.Get("userId")
		switch v.(type) {
		case int64:
			if v == 0 {
				//Session is invalid
				c.Redirect(http.StatusMovedPermanently, "/")
			}
		case nil:
			//User already logout
			c.Redirect(http.StatusMovedPermanently, "/")
		}

		// before request
		c.Next()

		// after request
		// ...
	}
}
