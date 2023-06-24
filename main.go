package main

import (
	"net/http"

	"github.com/J-khol-R/Email-Indexer/controllers"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/word/{keyWord}", controllers.GetEmails)

	handler := corsMiddleware(r)

	http.ListenAndServe(":8080", handler)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir todas las origenes
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Permitir los métodos HTTP especificados
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

		// Permitir los encabezados HTTP especificados
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Permitir el envío de cookies
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Si la solicitud es una solicitud OPTIONS, simplemente respondemos con los encabezados CORS sin continuar con la cadena de middleware
		if r.Method == "OPTIONS" {
			return
		}

		// Continuar con la cadena de middleware
		next.ServeHTTP(w, r)
	})
}
