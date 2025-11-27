package handlers

import (
	"encoding/json"
	"net/http"
	"student_management_service/database"
	"student_management_service/models"

	"github.com/go-chi/chi/v5"
)

// ==================== Program Handlers ====================

// GetPrograms retrieves all programs
func GetPrograms(w http.ResponseWriter, r *http.Request) {
	var programs []models.Program
	if err := database.DB.Preload("College").Find(&programs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(programs)
}

// GetProgramById retrieves a program by ID
func GetProgramById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var program models.Program
	if err := database.DB.Preload("College").First(&program, id).Error; err != nil {
		http.Error(w, "Program not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(program)
}

// CreateProgram creates a new program
func CreateProgram(w http.ResponseWriter, r *http.Request) {
	var program models.Program
	if err := json.NewDecoder(r.Body).Decode(&program); err != nil {
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

// UpdateProgram updates a program
func UpdateProgram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var program models.Program
	if err := database.DB.First(&program, id).Error; err != nil {
		http.Error(w, "Program not found", http.StatusNotFound)
		return
	}

	var updates models.Program
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Model(&program).Updates(updates).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload the program with College preloaded
	if err := database.DB.Preload("College").First(&program, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(program)
}

// DeleteProgram deletes a program
func DeleteProgram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var program models.Program
	if err := database.DB.First(&program, id).Error; err != nil {
		http.Error(w, "Program not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Delete(&program).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Program deleted successfully"})
}

// ==================== College Handlers ====================

// GetColleges retrieves all colleges
func GetColleges(w http.ResponseWriter, r *http.Request) {
	var colleges []models.College
	if err := database.DB.Find(&colleges).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(colleges)
}

// GetCollegeById retrieves a college by ID
func GetCollegeById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var college models.College
	if err := database.DB.First(&college, id).Error; err != nil {
		http.Error(w, "College not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(college)
}

// CreateCollege creates a new college
func CreateCollege(w http.ResponseWriter, r *http.Request) {
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

// UpdateCollege updates a college
func UpdateCollege(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var college models.College
	if err := database.DB.First(&college, id).Error; err != nil {
		http.Error(w, "College not found", http.StatusNotFound)
		return
	}

	var updates models.College
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Model(&college).Updates(updates).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload the college to get updated data
	if err := database.DB.First(&college, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(college)
}

// DeleteCollege deletes a college
func DeleteCollege(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var college models.College
	if err := database.DB.First(&college, id).Error; err != nil {
		http.Error(w, "College not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Delete(&college).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "College deleted successfully"})
}
