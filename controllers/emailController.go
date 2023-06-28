package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/J-khol-R/Email-Indexer/services"
	"github.com/go-chi/chi"
)

func GetEmails(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "keyWord")

	var datos map[string]int
	err := json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		http.Error(w, "Error al decodificar la respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := services.RequestZincsearch(key, datos["inicio"], datos["fin"])
	if err != nil {
		http.Error(w, "Error al traer datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
