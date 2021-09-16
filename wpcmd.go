package main

import (
	"embed"
	"github.com/samvaughton/wpcommand/v2/pkg/server"
)

//go:embed app/public
var staticFiles embed.FS

//go:embed config.yaml
var config string

//go:embed casbin/model.conf
var casbin string

func main() {
	server.Start(&staticFiles, config, casbin)
}
