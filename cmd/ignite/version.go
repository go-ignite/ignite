package main

import "github.com/fatih/color"

var Version string

var logo = `
┌─┐┌─┐  ┬┌─┐┌┐┌┬┌┬┐┌─┐
│ ┬│ │  ││ ┬││││ │ ├┤ 
└─┘└─┘  ┴└─┘┘└┘┴ ┴ └─┘

https://github.com/go-ignite/ignite
V%s

`

func displayVersion() {
	color.Cyan(logo, Version)
}
