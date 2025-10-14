package handlers

import (
	"encoding/json"
	"monolith/database"
	"monolith/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateProgram(w http.ResponseWriter, r *http.Request){
	var program models.Program

	if err := json.NewDecoder(r.Body).Decode(&program); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	if err := database.DB.Create(&program).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(program)
}

func GetProgram(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    var program models.Bank
    if err := database.DB.First(&program, id).Error; err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(program)
}