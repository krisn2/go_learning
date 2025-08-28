package main

import (
	"fmt"
	"log"

	"github.com/krisn2/ginserver/internals/auth"
	"github.com/krisn2/ginserver/internals/db"
	"github.com/krisn2/ginserver/internals/handlers"
	"github.com/krisn2/ginserver/internals/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	LoadConfig("config.json")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		AppConfig.Database.Host,
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.DBName,
		AppConfig.Database.Port,
		AppConfig.Database.SSLMode,
	)
	db.Connect(dsn)
	if err := db.DB.AutoMigrate(&models.User{}, &models.Item{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// CORS (adjust origins as needed)
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Authorization", "Content-Type"},
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	api := r.Group("/api")

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
	}

	userGroup := api.Group("/users")
	userGroup.Use(auth.RequireAuth())
	{
		userGroup.GET("/me", handlers.Me)
	}

	itemGroup := api.Group("/items")
	itemGroup.Use(auth.RequireAuth())
	{
		itemGroup.POST("", handlers.CreateItem)
		itemGroup.GET("/:id", handlers.GetItem)
		itemGroup.GET("", handlers.ListItems)
	}

	r.Run(":" + AppConfig.Server.Port)
}
