package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sabarish-manoharan/emp-management/db"
	"github.com/sabarish-manoharan/emp-management/routes"
	"github.com/spf13/viper"
)

func main() {
		loadEnv();
	port := viper.GetString("PORT")

	// Initialize router
	r := mux.NewRouter()

	// Connect to database
	db.ConnectDB()

	// Simple test route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	})

	// Register routes
	routes.RegisterEmployeeRoutes(r)
	routes.RegisterUserRoutes(r)
	routes.GetUserRoutes(r)
	routes.LoginUserRoutes(r)

	// Set up CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	// Start server with CORS
	fmt.Printf("Server running on port : %v\n", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatal(err)
	}
}

// func viperConfigFile(key string) string {
// 	viper.SetConfigFile(".env")
// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Fatalf("Error reading config file: %v", err)
// 	}
// 	value, ok := viper.Get(key).(string)
// 	if !ok {
// 		log.Fatalf("Invalid type assertion for key: %s", key)
// 	}
// 	return value
// }

func loadEnv() {
	viper.SetConfigFile(".env") // Try to read from .env file
	_ = viper.ReadInConfig()    // Ignore error if .env not found (Render won't have it)
	viper.AutomaticEnv()        // Read from OS env (Render injects these)
}
