package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/go-ignite/ignite"
)

var (
	versionFlag = flag.Bool("v", false, "version")
	version     = "unknown"
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		return
	}

	DisplayVersion()
	app, err := ignite.Init()
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	log.Fatal(app.Start())
}
