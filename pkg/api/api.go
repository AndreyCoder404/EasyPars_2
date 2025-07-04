package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupRouter configures and returns the Gin router with all API endpoints
// This function sets up the main router for the REST API
func SetupRouter() *gin.Engine {
	// Create Gin router with default middleware (Logger and Recovery)
	router := gin.Default()

	// Enable CORS for frontend integration
	// Future steps: Configure CORS properly for production
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API route group
	// Future steps: Add versioning (v1, v2), authentication middleware
	api := router.Group("/api")
	{
		// Health check endpoint
		// Future steps: Add database health check, system status
		api.GET("/health", handleHealth)

		// Fights endpoint - main functionality
		// Future steps: Add pagination, filtering, search capabilities
		api.GET("/fights", handleGetFights)

		// Future endpoints to be added:
		// api.GET("/fights/:id", handleGetFight)      // Get single fight
		// api.POST("/fights", handleCreateFight)      // Create new fight (admin)
		// api.PUT("/fights/:id", handleUpdateFight)   // Update fight (admin)
		// api.DELETE("/fights/:id", handleDeleteFight) // Delete fight (admin)
		// api.GET("/fighters", handleGetFighters)     // Get all fighters
		// api.GET("/fighters/:id", handleGetFighter)  // Get single fighter
	}

	// Serve static files for frontend
	// Future steps: Use proper static file server in production
	router.Static("/static", "./frontend")
	router.StaticFile("/", "./frontend/index.html")

	return router
}

// handleHealth handles GET requests to /api/health
// Returns the health status of the application
func handleHealth(c *gin.Context) {
	// Future steps: Add database connectivity check, parser status
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "EasyPars API is running",
		"version": "1.0.0",
	})
}

// handleGetFights handles GET requests for fight data
// Returns a list of fights (currently hardcoded)
func handleGetFights(c *gin.Context) {
	// Future steps:
	// 1. Integrate with parser to get real data
	// 2. Add database queries using GORM
	// 3. Implement pagination with query parameters
	// 4. Add filtering by date, fighter, location
	// 5. Add caching for performance

	// Mock data for initial testing
	fights := []map[string]interface{}{
		{
			"id":       1,
			"date":     "2024-01-15",
			"fighter1": "John Doe",
			"fighter2": "Jane Smith",
			"result":   "John Doe wins by KO",
			"location": "Las Vegas, NV",
			"round":    3,
			"time":     "2:45",
		},
		{
			"id":       2,
			"date":     "2024-01-20",
			"fighter1": "Mike Johnson",
			"fighter2": "Sarah Connor",
			"result":   "Sarah Connor wins by Decision",
			"location": "New York, NY",
			"round":    5,
			"time":     "5:00",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List of fights retrieved successfully",
		"data":    fights,
		"count":   len(fights),
	})
}

// Future functions to be implemented:
// - JWT authentication middleware
// - Input validation functions
// - Error handling middleware
// - Rate limiting middleware
// - Logging middleware
// - Database connection functions
