package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/Abiggj/pawwp/internal/account"
	"github.com/Abiggj/pawwp/internal/database"
	"github.com/Abiggj/pawwp/internal/middleware"
	"github.com/Abiggj/pawwp/internal/pet"
)

func main() {

	_ = godotenv.Load()

	database.Connect()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Dependency Injection
	accountRepo := account.NewRepository()
	accountService := account.NewService(accountRepo)
	accountHandler := account.NewHandler(accountService)

	// Public routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pawpp backend running 🐾"))
	})

	http.HandleFunc("/accounts/register", accountHandler.Register)
	http.HandleFunc("/accounts/login", accountHandler.Login)

// Protected route example
protected := middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	accountID := r.Context().Value(middleware.AccountIDKey)
	petRepo := pet.NewRepository()
	petService := pet.NewService(petRepo)
	petHandler := pet.NewHandler(petService)

	http.Handle("/pets", middleware.JWTAuth(http.HandlerFunc(petHandler.Create)))

	w.Write([]byte("Authenticated account ID: " + accountID.(string)))
}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){log.Println("PATH HIT:", r.URL.Path)})

	http.Handle("/protected", protected)
	log.Printf("🚀 Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
