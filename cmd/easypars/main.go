package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"easypars/pkg/api"
	"easypars/pkg/config"
	"github.com/gin-gonic/gin"
)

// Main entry point of the application
// This function initializes the application, loads configuration, and starts the server
func main() {
	// Initialize application logging
	log.Println("Starting EasyPars application...")

	// Load application configuration using Viper
	// This reads config.yaml and sets up all application settings
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Log successful configuration loading
	log.Printf("Configuration loaded successfully")
	log.Printf("Server will start on port: %s", cfg.Server.Port)

	// Initialize parser (future implementation)
	// Future steps: Initialize parser with configuration and setup goroutines
	// parser := parser.New(cfg.Parser)
	// parser.Start()
	log.Println("Parser initialization - pending implementation")

	// Initialize database connection (future implementation)
	// Future steps: Connect to PostgreSQL using GORM with loaded database config
	// db, err := database.Connect(cfg.Database)
	// if err != nil {
	//     log.Fatal("Failed to connect to database:", err)
	// }
	// defer db.Close()
	log.Println("Database connection - pending implementation")

	// Initialize API server with loaded configuration
	// This sets up all REST API endpoints using the Gin framework
	router := api.SetupRouter()

	// Configure Gin mode based on environment
	// Future steps: Add environment-specific configuration
	// if cfg.Environment == "production" {
	//     gin.SetMode(gin.ReleaseMode)
	// }

	// Add graceful shutdown handling
	// This ensures the application shuts down cleanly when receiving termination signals
	setupGracefulShutdown(router, cfg)

	// Prepare server address using the configured port
	// Ensures the port format is correct (adds : if not present)
	serverAddr := cfg.Server.Port
	if serverAddr[0] != ':' {
		serverAddr = ":" + serverAddr
	}

	// Start the HTTP server
	// Future steps: Add TLS support, custom timeouts, and middleware
	log.Printf("Server starting on port %s...", cfg.Server.Port)
	log.Printf("API endpoints available at: http://localhost%s/api/", serverAddr)
	log.Printf("Web interface available at: http://localhost%s/", serverAddr)

	// Start the server with the configured port
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupGracefulShutdown configures graceful shutdown for the application
// This function handles OS signals and ensures clean application termination
func setupGracefulShutdown(router *gin.Engine, cfg *config.Config) {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)

	// Register the channel to receive specific signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown signals
	go func() {
		// Wait for a signal
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)

		// Perform cleanup operations
		log.Println("Performing graceful shutdown...")

		// Future cleanup steps:
		// - Close database connections
		// - Stop parser goroutines
		// - Flush logs
		// - Save application state
		// - Close Redis connections
		// - Clean up temporary files

		log.Println("Shutdown complete")
		os.Exit(0)
	}()
}

// Future functions to be implemented:
// - initializeDatabase(cfg *config.Config) (*gorm.DB, error)
// - initializeParser(cfg *config.Config) (*parser.Parser, error)
// - initializeRedis(cfg *config.Config) (*redis.Client, error)
// - setupMiddleware(router *gin.Engine, cfg *config.Config)
// - setupRoutes(router *gin.Engine, cfg *config.Config)
// - initializeLogging(cfg *config.Config) error
// - validateEnvironment(cfg *config.Config) error
