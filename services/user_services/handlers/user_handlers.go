package handlers

import (
    "encoding/json"
    "net/http"
    "user_services/database"
    "user_services/models"
    "github.com/go-chi/chi/v5"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
    var users []models.UserData
    if err := database.DB.Find(&users).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func CreateUsers(w http.ResponseWriter, r *http.Request) {
    var users models.UserData
    if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := database.DB.Create(&users).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(users)
}

func GetuserById(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    var users models.UserData
    if err := database.DB.First(&users, id).Error; err != nil {
        http.Error(w, "Menu item not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}