package main

import (
	_ "chat-system/docs"
	"chat-system/internal/db"
	"chat-system/internal/routes"
	"chat-system/internal/utils"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Chat Microservice API
// @version 1.0
// @description This is a chat microservice.
// @host localhost:8080
// @BasePath /
func main() {
	utils.InitLogger()

	db.InitCassandra()
	db.InitRedis()

	router := routes.SetupRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Could not start server: %s\n", err.Error())
		utils.GetLogger().Error("Could not start server: ", err)
	}
}
