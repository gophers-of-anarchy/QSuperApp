package http

import (
	"QSuperApp/controllers"
	"QSuperApp/middlewares"
	"github.com/labstack/echo/v4"
	"log"
)

func StartServer(server *echo.Echo) {
	server.Use(middlewares.AuthMiddleware)
	server.GET("/admin/orders", controllers.DecideOrderStatusHandler)
	server.GET("/admin/orders/status", controllers.ChangeOrderStatusHandler)
	log.Fatal(server.Start("localhost:6060"))
}
