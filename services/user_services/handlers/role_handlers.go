package handlers

import (
	"encoding/json"
	"net/http"
	"user_services/database"
	"user_services/models"

	"github.com/go-chi/chi/v5"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	var roles []models.UserRole
	if err := database.DB.Find(&roles).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.UserRole
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&role).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func GetRoleById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var role models.UserRole
	if err := database.DB.First(&role, id).Error; err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var role models.UserRole
	if err := database.DB.First(&role, id).Error; err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Save(&role).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var role models.UserRole
	if err := database.DB.First(&role, id).Error; err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}