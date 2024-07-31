package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danilsgit/test-backend/routes"
	"github.com/go-chi/chi/v5"
)

func main() {

	// Create a new router
	r := chi.NewRouter()

	// create cors
	cors := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	}

	// Use the cors
	r.Use(cors)

	// Define the routes
	r.Get("/", routes.HomeHandler)
	r.Get("/search", routes.SearchHandler)
	r.Get("/getAll", routes.GetAllHandler)

	// Get the port
	PORT := os.Getenv("PORT")
	fmt.Println("Server running on port: ", PORT)

	// Start the server
	http.ListenAndServe(":"+PORT, r)

}
