package handlers

import (
	"encoding/json"
	"net/http"
	"student_management_service/database"
	"student_management_service/models"

	"github.com/go-chi/chi/v5"
)

// GetStudents retrieves all students
func GetStudents(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	if err := database.DB.Preload("Program").Preload("College").Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// CreateStudent creates a new student
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if student.RubIDCardNumber == "" || student.Name == "" || student.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&student).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload the student with Program and College preloaded
	if err := database.DB.Preload("Program").Preload("College").First(&student, student.ID).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// GetStudentById retrieves a student by database ID
func GetStudentById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var student models.Student
	if err := database.DB.Preload("Program").Preload("College").First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// GetStudentByRubId retrieves a student by their RUB ID card number
func GetStudentByRubId(w http.ResponseWriter, r *http.Request) {
	rubId := chi.URLParam(r, "rubId")

	var student models.Student
	if err := database.DB.Preload("Program").Preload("College").Where("rub_id_card_number = ?", rubId).First(&student).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// GetStudentsByProgram retrieves all students in a specific program
func GetStudentsByProgram(w http.ResponseWriter, r *http.Request) {
	programId := chi.URLParam(r, "programId")

	var students []models.Student
	if err := database.DB.Preload("Program").Preload("College").Where("program_id = ?", programId).Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetStudentsByCollege retrieves all students in a specific college
func GetStudentsByCollege(w http.ResponseWriter, r *http.Request) {
	collegeId := chi.URLParam(r, "collegeId")

	var students []models.Student
	if err := database.DB.Preload("Program").Preload("College").Where("college_id = ?", collegeId).Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// UpdateStudent updates a student by ID
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var student models.Student
	if err := database.DB.First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	var updates models.Student
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Model(&student).Updates(updates).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload the student with Program and College preloaded
	if err := database.DB.Preload("Program").Preload("College").First(&student, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// DeleteStudent soft deletes a student by ID
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var student models.Student
	if err := database.DB.First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Delete(&student).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Student deleted successfully"})
}

// SearchStudents searches students by name, email, or RUB ID card number
func SearchStudents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	var students []models.Student
	searchPattern := "%" + query + "%"
	if err := database.DB.Preload("Program").Preload("College").Where("name ILIKE ? OR email ILIKE ? OR rub_id_card_number ILIKE ?", 
		searchPattern, searchPattern, searchPattern).Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
