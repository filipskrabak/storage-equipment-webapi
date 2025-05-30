package main

import (
	"log"
	"os"
	"strings"

	"context"
	"time"

	"github.com/filipskrabak/storage-equipment-webapi/api"
	"github.com/filipskrabak/storage-equipment-webapi/internal/db_service"
	"github.com/filipskrabak/storage-equipment-webapi/internal/storage_equipment"
	"github.com/gin-contrib/cors"
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
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update  middleware
	equipmentDbService := db_service.NewMongoService[storage_equipment.EquipmentItem](db_service.MongoServiceConfig{
		Collection: "equipment"
	})
	defer equipmentDbService.Disconnect(context.Background())

	ordersDbService := db_service.NewMongoService[storage_equipment.Order](db_service.MongoServiceConfig{
		Collection: "orders",
	})
	defer ordersDbService.Disconnect(context.Background())

	engine.Use(func(ctx *gin.Context) {
		ctx.Set("equipment_db_service", equipmentDbService)
		ctx.Set("orders_db_service", ordersDbService)
		ctx.Next()
	})
	// request routings
	handleFunctions := &storage_equipment.ApiHandleFunctions{
		EquipmentManagementAPI: storage_equipment.NewEquipmentManagementApi(),
		EquipmentOrdersAPI:    storage_equipment.NewEquipmentOrdersApi(),
	}
	storage_equipment.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
