package handlers

import (
	"encoding/json"
	"monolith/database"
	"monolith/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateBank(w http.ResponseWriter, r *http.Request) {
	var bank models.Bank

	if err := json.NewDecoder(r.Body).Decode(&bank); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	if err := database.DB.Create(&bank).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(bank)
}

func GetBank(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    var bank models.Bank
    if err := database.DB.First(&bank, id).Error; err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bank)
}