package main

import (
	"log"
	"os"
	"strings"

	"github.com/filipskrabak/storage-equipment-webapi/api"
	"github.com/filipskrabak/storage-equipment-webapi/internal/storage_equipment"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	// request routings
	handleFunctions := &storage_equipment.ApiHandleFunctions{
		EquipmentManagementAPI: storage_equipment.NewEquipmentManagementApi(),
	}
	storage_equipment.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
