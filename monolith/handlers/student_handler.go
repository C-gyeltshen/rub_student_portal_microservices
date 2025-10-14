package handlers

import (
	"encoding/json"
	"monolith/database"
	"monolith/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)


func CreateStudent(w http.ResponseWriter, r *http.Request){
	var student models.Student

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	if err := database.DB.Create(&student).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(student)
}

func GetStudent(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    var student models.Student
    if err := database.DB.First(&student, id).Error; err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(student)
}