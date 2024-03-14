package main

import (
	"LarsWebV0/dao"
	"LarsWebV0/router"
	"context"

	logger "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/*func main() {
	fmt.Print("hello Lars!")
}*/

func init() {
	dao.SetupDB()
	//dao.MongoDBSetup()
	dao.EsSetup()
}

// @title LarsWebv0
// @version 1.0
// @description 小黄成长日记

// @contact.name Laurus
// @contact.email 2312593392@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	router := router.SetupRouter()

	//--------------------------------------something about shutdown-----------------------------------

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	/*	go func() {
		if err := httpServer.ListenAndServeTLS("weplant.top_bundle.pem", "weplant.top.key"); err != nil && err != http.ErrServerClosed {
			logger.Errorf("listen: %s\n", err)
			return
		}
	}()*/

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("listen: %s\n", err)
			return
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Infof("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		dao.CloseDB()
	}()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("Server Shutdown: %v", err)
		return
	}
	logger.Infof("Server exiting")
}
