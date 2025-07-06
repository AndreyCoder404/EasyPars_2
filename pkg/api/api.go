package api

import (
	"easypars/pkg/parser"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

		// Fights endpoint - main functionality with real parser integration
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

// handleGetFights handles GET requests for fight data using the real parser
// Integrates with the parser to fetch live data from vringe.com
func handleGetFights(c *gin.Context) {
	log.Println("Received request for fight data, initializing parser")

	// Create a new parser instance with the target URL
	// Using hardcoded URL as specified in requirements
	parserInstance := parser.NewParser("https://vringe.com/results/")

	// Call ParseFights to get real data from the website
	// This performs HTTP request, HTML parsing, and concurrent processing
	fights, err := parserInstance.ParseFights()
	if err != nil {
		// Log the error for debugging
		log.Printf("Error parsing fights: %v", err)

		// Return appropriate HTTP error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse fight data",
			"message": "Unable to retrieve fight information at this time",
			"details": err.Error(),
		})
		return
	}

	// Check if we got any results
	if len(fights) == 0 {
		log.Println("No fight data found, returning empty result")
		c.JSON(http.StatusOK, gin.H{
			"message": "No fight data available",
			"data":    []interface{}{},
			"count":   0,
		})
		return
	}

	// Log successful response
	log.Printf("Successfully retrieved %d fights from parser", len(fights))

	// Return successful response with parsed fight data
	c.JSON(http.StatusOK, gin.H{
		"message": "Fight data retrieved successfully from vringe.com",
		"data":    fights,
		"count":   len(fights),
		"source":  "vringe.com/results/",
	})

	// Future steps:
	// 1. Add caching to reduce load on target website
	// 2. Implement pagination with query parameters (?page=1&limit=10)
	// 3. Add filtering by date range (?from=2024-01-01&to=2024-12-31)
	// 4. Add search functionality (?search=fighter_name)
	// 5. Add sorting options (?sort=date&order=desc)
	// 6. Store parsed data in database for faster subsequent requests
	// 7. Add rate limiting to prevent abuse
	// 8. Add request timeout handling
	// 9. Add data validation and sanitization
	// 10. Add metrics and monitoring for parsing performance
}
