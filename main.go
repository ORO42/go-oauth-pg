package main

import (
	"fmt"
	"gop/auth"
	"gop/db"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	defer db.CloseDB()
	auth.InitAuth()

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/auth/google", auth.BeginAuthController)
	mux.HandleFunc("/auth/google/callback", auth.CompleteAuthController)
	mux.HandleFunc("/logout", auth.LogoutController)

	// Start server
	fmt.Println("Starting server on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}

// Simple root route
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the auth server")
	fmt.Fprintln(w, "Visit /auth/google to login")
}
