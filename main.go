package main

import (
	"QSuperApp/configs"
	"QSuperApp/controllers"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Check application environment
	if os.Getenv("ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("App .env file not found")
		}
	}

	// Connect postgres
	configs.ConnectToDatabase()

	e := echo.New()
	apiGroup := e.Group("/api/v1")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Airplane routes
	airplaneApiGroup := apiGroup.Group("/airplane")
	airplaneApiGroup.POST("/add", controllers.Add)
	airplaneApiGroup.PUT("/update", controllers.Update)
	airplaneApiGroup.GET("/all", controllers.GetAll)
	airplaneApiGroup.GET("/:id", controllers.Get)
	airplaneApiGroup.DELETE("/:id", controllers.Delete)

	// Run Server
	e.Logger.Fatal(e.Start(":8080"))
}
