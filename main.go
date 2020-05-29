package main

import (
	"iMonitor/conf"
	_ "iMonitor/docs"
	"iMonitor/router"
	"os"
)

func main() {
	conf.Init()

	r := router.InitRouter()
	r.Run(os.Getenv("PORT"))
}
