package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ValidateSession() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		v := session.Get("userId")
		fmt.Println("In middleware handler...")

		switch v.(type) {
		case int64:
			fmt.Println("session --> userId is:", v)
			if v == 0 {
				//Session is invalid
				c.Redirect(http.StatusFound, "/")
			}
			c.Set("userId", v.(int64))
		case nil:
			//User already logout
			fmt.Println("session --> empty session")
			c.Redirect(http.StatusFound, "/")
		default:
			fmt.Println("session --> unknown session")
			c.Redirect(http.StatusFound, "/")
		}

		// before request
		c.Next()

		// after request
		// ...
	}
}
