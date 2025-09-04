package main

import (
	"context"
	"fmt"
	"time"

	"github.com/21keshav/IBackendApplication/config"
	"github.com/21keshav/IBackendApplication/controller"
	"github.com/21keshav/IBackendApplication/resources/bidManager"
	"github.com/21keshav/IBackendApplication/resources/project"
	"github.com/21keshav/IBackendApplication/util"
	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	glog.Info("main")
	glog.InfoDepth(1, "started")
	defer glog.InfoDepth(1, "completed")
	e := echo.New()
	var conf config.Config
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		glog.Error("Error reading config file")
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	mongoURL := fmt.Sprintf("%s:%s", conf.Database.Server, conf.Database.Port)
	mongoClient := util.NewMongoClient(context.TODO(), mongoURL)
	context, _ := context.WithTimeout(context.Background(), 60*time.Second)
	projectManager := project.NewProjectManager(mongoClient, context, conf.DatabaseDetails)
	bidManager := bidManager.NewBidManager(projectManager, context)
	controller := controller.NewController(bidManager, projectManager)
	controller.AttachHandlers(e)

	e.Start(":1234")
}
