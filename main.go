package main

import (
	"fmt"
	"log"
	"net/http"
	"rest-services/config"
	"rest-services/controllers"
	"rest-services/db"
	"rest-services/middleware"
)

func main() {
	config.Initialize()

	// Migrate the database
	db.Migrate()

	// Rroutes
	http.HandleFunc("/signup", controllers.SignUp)
	http.HandleFunc("/signin", controllers.SignIn)
	http.HandleFunc("/refresh", controllers.RefreshToken)

	// Protected route with middleware
	http.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(controllers.ProtectedEndpoint)))

	// Revoke token endpoint
	http.HandleFunc("/revoke", controllers.RevokeToken)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
