package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/robstave/go-template/docs"
	"github.com/robstave/go-template/internal/adapters/controller"
	"github.com/robstave/go-template/internal/adapters/repositories"
	"github.com/robstave/go-template/internal/domain"
	"github.com/robstave/go-template/internal/domain/types"
	"github.com/robstave/go-template/internal/logger"
	httpSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title go-template
// @version 1.0
// @description API documentation for go-template
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize Logger
	slogger := logger.InitializeLogger()
	logger.SetLogger(slogger) // Optional: If you prefer setting a package-level logger

	// Initialize Database
	dbPath := "./go-template.db"
	if path := os.Getenv("DB_PATH"); path != "" {
		dbPath = path
	}
	slogger.Info("DBPath set", "dbpath", dbPath)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		slogger.Error("Failed to connect to database", "path", dbPath, "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&types.Widget{}) // Add this line
	if err != nil {
		slogger.Error("Failed to migrate database", "error", err)
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Repositories
	repo := repositories.NewRepositorySQLite(db)

	// Initialize Service
	service := domain.NewService(slogger, repo)

	// Initialize Controller
	controller := controller.NewController(service, slogger)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())
	/*
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*", "http://localhost:8789", "http://192.168.86.176:8789", "http://localhost:3000"}, // Update as needed
			AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
			AllowHeaders: []string{
				"Content-Type",
				"Authorization",
			},
		}))
	*/

	// API Routes
	api := e.Group("/api")
	widgetGroup := api.Group("/widgets")

	// Widget routes
	widgetGroup.POST("", controller.CreateWidget)
	widgetGroup.GET("", controller.GetAllWidgets)
	widgetGroup.GET("/:id", controller.GetWidget)
	widgetGroup.PUT("/:id", controller.UpdateWidget)
	widgetGroup.DELETE("/:id", controller.DeleteWidget)

	// Swagger endpoint
	e.GET("/swagger/*", httpSwagger.WrapHandler)

	// Catch-all route to serve index.html for client-side routing
	e.GET("/*", func(c echo.Context) error {
		slogger.Info("other url", "path", c.Request().URL.Path)
		// Prevent Echo from serving API routes with this handler
		if strings.HasPrefix(c.Request().URL.Path, "/api/") {
			return c.NoContent(http.StatusNotFound)
		}
		path := c.Request().URL.Path

		// Redirect /swagger to /swagger/index.html#/
		if path == "/swagger" || path == "/swagger/" {
			return c.Redirect(http.StatusFound, "/swagger/index.html#/")
		}

		return c.NoContent(http.StatusNotFound)
	})

	// Start Server
	port := "8711"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	slogger.Info("Starting server", "port", port)
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		slogger.Error("Shutting down the server", "error", err)
		log.Fatalf("Shutting down the server: %v", err)
	}
}
