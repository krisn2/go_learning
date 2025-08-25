package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/krisn2/go-social/config"
	"github.com/krisn2/go-social/handlers"
	"github.com/krisn2/go-social/middleware"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db, cfg)
	postHandler := handlers.NewPostHandler(db)
	likeHandler := handlers.NewLikeHandler(db)
	commentHandler := handlers.NewCommentHandler(db)

	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public user routes
		users := api.Group("/users")
		{
			users.GET("/", userHandler.ListUsers) // Admin gated
		}

		// Public post routes
		posts := api.Group("/posts")
		{
			posts.GET("/", postHandler.List)
			posts.GET("/:id", postHandler.Get)
			posts.GET("/:id/comments", commentHandler.List)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			// User routes
			protectedUsers := protected.Group("/users")
			{
				protectedUsers.GET("/me", userHandler.GetMe)
				protectedUsers.PATCH("/me", userHandler.UpdateMe)
				protectedUsers.DELETE("/me", userHandler.DeleteMe)
			}

			// Post routes
			protectedPosts := protected.Group("/posts")
			{
				protectedPosts.POST("/", postHandler.Create)
				protectedPosts.PATCH("/:id", postHandler.Update)
				protectedPosts.DELETE("/:id", postHandler.Delete)
				protectedPosts.POST("/:id/like", likeHandler.Toggle)
				protectedPosts.POST("/:id/comments", commentHandler.Create)
			}

			// Comment routes
			protectedComments := protected.Group("/comments")
			{
				protectedComments.DELETE("/:id", commentHandler.Delete)
			}
		}
	}

	return router
}
