package main

import (
	"QSuperApp/configs"
	"QSuperApp/controllers"
	"QSuperApp/middlewares"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

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

	// Load templates
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Auth routes
	auth := apiGroup.Group("/auth")
	auth.POST("/register", controllers.RegisterHandler)
	auth.POST("/login", controllers.LoginHandler)

	// Airplane routes
	airplaneApiGroup := apiGroup.Group("/airplane")
	airplaneApiGroup.POST("/add", controllers.Add)
	airplaneApiGroup.PUT("/update", controllers.Update)
	airplaneApiGroup.GET("/all", controllers.GetAll)
	airplaneApiGroup.GET("/:id", controllers.Get)
	airplaneApiGroup.DELETE("/:id", controllers.Delete)

	// Order routes
	orderManagementApiGroup := apiGroup.Group("/order-management")
	orderManagementApiGroup.POST("/admin/orders", controllers.DecideOrderStatusHandler)
	orderManagementApiGroup.POST("/admin/orders/status", controllers.ChangeOrderStatusHandler)
	orderManagementApiGroup.GET("/admin/orders/list", controllers.GetAllOrderHandler)

	// Payment routes
	paymentApiGroup := apiGroup.Group("/payment")
	paymentApiGroup.Use(middlewares.AuthMiddleware)
	paymentApiGroup.POST("/advance", controllers.AdvancePaymentHandler)
	paymentApiGroup.POST("/finalize", controllers.FinalPaymentHandler)
	paymentApiGroup.GET("/orders/:orderId", controllers.VerifyPaymentHandler)

	verifyPaymentGroup := apiGroup.Group("/verify")
	verifyPaymentGroup.POST("/page", controllers.VerifyPaymentPageHandler)
	verifyPaymentGroup.POST("/payment", controllers.VerifyPaymentHandler)

	// Run Server
	e.Logger.Fatal(e.Start(":8080"))
}
