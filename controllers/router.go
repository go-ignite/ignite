package controllers

import "github.com/gin-gonic/gin"

type MainRouter struct {
	router *gin.Engine
}

func (self *MainRouter) InitRouter(r *gin.Engine) {
	self.router = r

	self.router.GET("/", self.IndexHandler)
	self.router.Run(":5000")
}
