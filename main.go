package main

import (
	"fmt"
	"ignite/controllers"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	workingDir string
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	pwd, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Working directory:", pwd)
	workingDir = pwd

	mainRouter := &controllers.MainRouter{}
	mainRouter.Initialize(r)
}
