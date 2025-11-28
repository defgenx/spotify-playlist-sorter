package api

import (
	"github.com/gin-gonic/gin"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/api/handlers"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/api/middleware"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/config"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/service"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/session"
	spotifyClient "github.com/adelvecchio/spotify-playlist-sorter/internal/spotify"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/sse"
)

// Router sets up and returns the Gin router
func NewRouter(
	cfg *config.Config,
	spotifyClient *spotifyClient.Client,
	sessionStore *session.Store,
	broadcaster *sse.Broadcaster,
	libraryService *service.LibraryService,
	sorterService *service.SorterService,
	executorService *service.ExecutorService,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware(cfg.Server.AllowOrigins))

	// Create handlers
	authHandler := handlers.NewAuthHandler(spotifyClient, sessionStore, cfg)
	libraryHandler := handlers.NewLibraryHandler(spotifyClient, libraryService)
	sortHandler := handlers.NewSortHandler(spotifyClient, libraryService, sorterService, executorService)
	eventsHandler := handlers.NewEventsHandler(broadcaster)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.GET("/login", authHandler.Login)
			auth.GET("/callback", authHandler.Callback)
			auth.GET("/complete", authHandler.CompleteLogin)
			auth.POST("/logout", middleware.AuthMiddleware(sessionStore), authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(sessionStore), authHandler.GetMe)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(sessionStore))
		{
			// Library routes
			library := protected.Group("/library")
			{
				library.GET("/analysis", libraryHandler.GetAnalysis)
			}

			// Sort routes
			sort := protected.Group("/sort")
			{
				sort.POST("/plan", sortHandler.GeneratePlan)
				sort.POST("/execute", sortHandler.ExecutePlan)
			}

			// Events routes (SSE)
			events := protected.Group("/events")
			{
				events.GET("", eventsHandler.StreamEvents)
				events.POST("/test", eventsHandler.TestEvent)
			}
		}
	}

	return router
}
