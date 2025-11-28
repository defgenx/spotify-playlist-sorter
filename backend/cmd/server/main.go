package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/adelvecchio/spotify-playlist-sorter/internal/api"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/config"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/service"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/session"
	spotifyClient "github.com/adelvecchio/spotify-playlist-sorter/internal/spotify"
	"github.com/adelvecchio/spotify-playlist-sorter/internal/sse"
)

func main() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	log.Info().Msg("Starting Spotify Playlist Sorter API")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	log.Info().
		Str("port", cfg.Server.Port).
		Str("frontendURL", cfg.Server.FrontendURL).
		Msg("Configuration loaded")

	// Initialize Spotify client
	spotifyClient := spotifyClient.NewClient(
		cfg.Spotify.ClientID,
		cfg.Spotify.ClientSecret,
		cfg.Spotify.RedirectURL,
	)
	log.Info().Msg("Spotify client initialized")

	// Initialize session store
	sessionStore := session.NewStore()
	log.Info().Msg("Session store initialized")

	// Initialize SSE broadcaster
	broadcaster := sse.NewBroadcaster()
	log.Info().Msg("SSE broadcaster initialized")

	// Initialize services
	libraryService := service.NewLibraryService(spotifyClient, broadcaster)
	sorterService := service.NewSorterService(libraryService)
	executorService := service.NewExecutorService(spotifyClient, libraryService, broadcaster)
	log.Info().Msg("Services initialized")

	// Create router
	router := api.NewRouter(
		cfg,
		spotifyClient,
		sessionStore,
		broadcaster,
		libraryService,
		sorterService,
		executorService,
	)
	log.Info().Msg("Router configured")

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Minute, // Long timeout for SSE
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Info().Str("addr", srv.Addr).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
