package handlers

import (
	// "encoding/json"
	"encoding/json"
	"monolith/database"
	"monolith/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateCollege(w http.ResponseWriter, r*http.Request){
     var college models.College

	if err := json.NewDecoder(r.Body).Decode(&college); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := database.DB.Create(&college).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(college)
}

func GetCollege(w http.ResponseWriter, r *http.Request){
    id := chi.URLParam(r, "id")

    var college models.College
    if err := database.DB.First(&college, id).Error; err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(college)
}