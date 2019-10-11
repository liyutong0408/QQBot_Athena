package main

import (
	"Athena/conf"
	"Athena/server"
)

func main() {
	conf.Init()

	r := server.NewRouter()

	r.Run(":65321")
}
