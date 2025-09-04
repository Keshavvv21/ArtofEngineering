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
	// ---- Application Startup Logs ----
	glog.Info("Starting IBackendApplication...")
	defer glog.Info("Application stopped.")

	// ---- Load Configuration ----
	var conf config.Config
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		glog.Error("Failed to read config file: ", err)
		// TODO: consider panic/exit if config is critical
	}

	// ---- Initialize Echo Web Framework ----
	e := echo.New()
	e.Use(middleware.Logger())   // Log all HTTP requests
	e.Use(middleware.Recover())  // Recover from panics and return HTTP 500

	// ---- Setup Database Connection (MongoDB) ----
	mongoURL := fmt.Sprintf("%s:%s", conf.Database.Server, conf.Database.Port)
	mongoClient := util.NewMongoClient(context.Background(), mongoURL)

	// Create a context with timeout for DB operations
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// ---- Initialize Resource Managers ----
	// Project Manager handles project-related operations
	projectManager := project.NewProjectManager(mongoClient, ctx, conf.DatabaseDetails)

	// Bid Manager handles bidding logic, depends on ProjectManager
	bidManager := bidManager.NewBidManager(projectManager, ctx)

	// ---- Setup Controller & Route Handlers ----
	// Controller wires HTTP routes to application logic
	ctrl := controller.NewController(bidManager, projectManager)
	ctrl.AttachHandlers(e)

	// ---- Start HTTP Server ----
	// TODO: replace with graceful shutdown (e.Shutdown) for production use
	port := ":1234"
	glog.Infof("Server listening on %s", port)
	if err := e.Start(port); err != nil {
		glog.Errorf("Error starting server: %v", err)
	}
}
