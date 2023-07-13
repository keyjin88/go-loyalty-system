package main

import (
	"github.com/keyjin88/go-loyalty-system/internal/app"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
)

func main() {
	server := app.New()
	// api server start
	err := server.Start()
	if nil != err {
		logger.Log.Info("Error starting api server")
		return
	}
}
