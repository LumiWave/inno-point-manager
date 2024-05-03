package main

import (
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Error(err)
		return
	}

	if err := app.Start(); err != nil {
		log.Error(err)
		return
	}
}
