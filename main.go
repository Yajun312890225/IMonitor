package main

import (
	"iMonitor/conf"
	_ "iMonitor/docs"
	"iMonitor/router"
)

func main() {
	conf.Init()

	r := router.InitRouter()
	r.Run(":9528")
}
