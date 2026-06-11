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
	"github.com/Abiggj/pawwp/internal/post"
)

func main() {
	_ = godotenv.Load()

	database.Connect()

	// Auto Migration
	database.DB.AutoMigrate(&account.Account{}, &pet.Pet{}, &post.Post{}, &post.Boop{}, &post.Woof{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Dependency Injection
	accountRepo := account.NewRepository()
	accountService := account.NewService(accountRepo)
	accountHandler := account.NewHandler(accountService)

	petRepo := pet.NewRepository()
	petService := pet.NewService(petRepo)
	petHandler := pet.NewHandler(petService)

	postRepo := post.NewRepository()
	postService := post.NewService(postRepo, petRepo)
	postHandler := post.NewHandler(postService)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pawpp backend running 🐾"))
	})

	mux.HandleFunc("POST /accounts/register", accountHandler.Register)
	mux.HandleFunc("POST /accounts/login", accountHandler.Login)
	mux.HandleFunc("POST /accounts/refresh", accountHandler.Refresh)

	// Protected Account routes
	mux.Handle("GET /accounts/me", middleware.JWTAuth(http.HandlerFunc(accountHandler.Me)))
	mux.Handle("PUT /accounts/me", middleware.JWTAuth(http.HandlerFunc(accountHandler.Update)))
	mux.Handle("DELETE /accounts/me", middleware.JWTAuth(http.HandlerFunc(accountHandler.Delete)))

	// Protected Pet routes
	mux.Handle("POST /pets", middleware.JWTAuth(http.HandlerFunc(petHandler.Create)))
	mux.Handle("GET /pets/{id}", middleware.JWTAuth(http.HandlerFunc(petHandler.Get)))
	mux.Handle("PUT /pets/{id}", middleware.JWTAuth(http.HandlerFunc(petHandler.Update)))
	mux.Handle("DELETE /pets/{id}", middleware.JWTAuth(http.HandlerFunc(petHandler.Delete)))
	mux.Handle("GET /accounts/{account_id}/pets", middleware.JWTAuth(http.HandlerFunc(petHandler.ListByAccount)))

	// Post routes
	mux.Handle("POST /posts", middleware.JWTAuth(http.HandlerFunc(postHandler.Create)))
	mux.Handle("GET /feed", middleware.JWTAuth(http.HandlerFunc(postHandler.GetFeed)))
	mux.Handle("GET /pets/{pet_id}/showcase", middleware.JWTAuth(http.HandlerFunc(postHandler.GetShowcase)))
	mux.Handle("PUT /posts/{post_id}/showcase", middleware.JWTAuth(http.HandlerFunc(postHandler.ToggleShowcase)))
	mux.Handle("GET /pets/{pet_id}/archive", middleware.JWTAuth(http.HandlerFunc(postHandler.GetArchive)))
	mux.Handle("POST /posts/{post_id}/boops", middleware.JWTAuth(http.HandlerFunc(postHandler.ToggleBoop)))
	mux.Handle("POST /posts/{post_id}/woofs", middleware.JWTAuth(http.HandlerFunc(postHandler.AddWoof)))
	mux.Handle("GET /posts/{post_id}/woofs", middleware.JWTAuth(http.HandlerFunc(postHandler.GetWoofs)))

	// Debug/Catch-all
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("PATH HIT:", r.URL.Path)
	})

	log.Printf("🚀 Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
