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
	if err := database.DB.Preload("Role").Find(&users).Error; err != nil {
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

	// Reload the user with the Role relationship
	database.DB.Preload("Role").First(&users, users.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(users)
}

func GetuserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var users models.UserData
	if err := database.DB.Preload("Role").First(&users, id).Error; err != nil {
		http.Error(w, "Menu item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUsersByRoleId(w http.ResponseWriter, r *http.Request) {
	roleId := chi.URLParam(r, "roleId")

	var users []models.UserData
	if err := database.DB.Preload("Role").Where("user_role_id = ?", roleId).Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func DeleteUsersByRoleId(w http.ResponseWriter, r *http.Request) {
	roleId := chi.URLParam(r, "roleId")

	// Delete all users with the specified role ID
	result := database.DB.Where("user_role_id = ?", roleId).Delete(&models.UserData{})

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":       "Users deleted successfully",
		"deleted_count": result.RowsAffected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateFinanceOfficer(w http.ResponseWriter, r *http.Request) {
	var user models.UserData
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set role ID to 2 (Finance officer)
	user.UserRoleID = 2

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload the user with the Role relationship
	database.DB.Preload("Role").First(&user, user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
