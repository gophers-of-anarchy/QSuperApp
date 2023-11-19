package main

import (
	"QSuperApp/configs"
	"QSuperApp/controllers"
	"QSuperApp/middlewares"
	services "QSuperApp/services/sms"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Server Starting...")

	// Check application environment
	if os.Getenv("ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("App .env file not found")
		}
	}

	// Connect postgres
	configs.ConnectToDatabase()

	// Connect redis
	configs.ConnectToRedis()

	// Register services
	services.RegisterSMSService()

	// Define mux router
	router := mux.NewRouter()

	// Register and login endpoints
	router.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", controllers.LoginHandler).Methods("POST")

	// Auth middleware
	router.Use(middlewares.AuthMiddleware)

	// Restricted area endpoint
	router.HandleFunc("/restricted_area", controllers.RestrictedAreaHandler).Methods("GET")

	// Card endpoints
	router.HandleFunc("/create_card", controllers.CreateCardHandler).Methods("POST")

	// Account endpoints
	router.HandleFunc("/create_account", controllers.CreateAccountHandler).Methods("POST")
	
	// Order endpoints
	router.HandleFunc("/create/order", controllers.CreateOrderHandler).Methods("POST")

	// Money transfer endpoint
	router.HandleFunc("/transfer", controllers.TransferHandler).Methods("POST")

	// Top transaction users endpoint
	router.HandleFunc("/top_users", controllers.TopUsersHandler).Methods("GET")

	log.Println("Server is Running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
