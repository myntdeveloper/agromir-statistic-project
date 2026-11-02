package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/myntdeveloper/agromir_boiler_server/internal/config"
	"github.com/myntdeveloper/agromir_boiler_server/internal/handler"
	"github.com/myntdeveloper/agromir_boiler_server/internal/middleware"
	"github.com/myntdeveloper/agromir_boiler_server/internal/repo"
	"github.com/myntdeveloper/agromir_boiler_server/internal/service"
)

func SetupRoutes(r *gin.Engine, conf *config.Config, boilerHandler *handler.BoilerHandler, ingestHandler *handler.IngestHandler) {
	api := r.Group("/api/v1")

	api.GET("/warehouses", boilerHandler.ListWarehouses)
	api.GET("/warehouses/:id/boilers", boilerHandler.ListBoilers)
	api.GET("/boilers/:id", boilerHandler.GetBoiler)
	api.GET("/boilers/:id/trend", boilerHandler.GetTrend)
	api.GET("/plants/:plant_id/status", ingestHandler.GetPlantStatus)

	ingest := api.Group("/ingest")
	ingest.Use(middleware.AgentBearer(conf.AgentAPIKey))
	ingest.POST("/readings", ingestHandler.IngestReadings)
}

func main() {
	conf := config.Load()
	db := repo.NewDB()
	defer db.Close()

	if err := repo.Migrate(db); err != nil {
		log.Fatal("Migration Error:", err)
	}

	boilerRepo := repo.NewBoilerRepo(db)
	boilerService := service.NewBoilerService(boilerRepo)
	boilerHandler := handler.NewBoilerHandler(boilerService)

	ingestRepo := repo.NewIngestRepo(db)
	ingestService := service.NewIngestService(ingestRepo)
	ingestHandler := handler.NewIngestHandler(ingestService)

	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		AllowCredentials: false,
	}))

	SetupRoutes(r, conf, boilerHandler, ingestHandler)

	log.Printf("Server started on :%s", conf.Port)
	if err := r.Run(":" + conf.Port); err != nil {
		log.Fatal(err)
	}
}
