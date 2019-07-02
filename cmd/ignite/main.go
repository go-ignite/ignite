package main

import (
	"flag"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/go-ignite/ignite"
)

var (
	versionFlag = flag.Bool("v", false, "version")
)

func main() {
	flag.Parse()
	displayVersion()

	if *versionFlag {
		return
	}

	app, err := ignite.Init()
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	log.Fatal(app.Start())
}
