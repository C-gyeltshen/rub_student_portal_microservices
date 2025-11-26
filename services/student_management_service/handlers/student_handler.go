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
	if err := database.DB.Find(&students).Error; err != nil {
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

	if student.StudentID == "" || student.FirstName == "" || student.LastName == "" || student.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
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

// GetStudentById retrieves a student by database ID
func GetStudentById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var student models.Student
	if err := database.DB.First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// GetStudentByStudentId retrieves a student by their student ID
func GetStudentByStudentId(w http.ResponseWriter, r *http.Request) {
	studentId := chi.URLParam(r, "studentId")

	var student models.Student
	if err := database.DB.Where("student_id = ?", studentId).First(&student).Error; err != nil {
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
	if err := database.DB.Where("program_id = ?", programId).Find(&students).Error; err != nil {
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
	if err := database.DB.Where("college_id = ?", collegeId).Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetStudentsByStatus retrieves all students with a specific status
func GetStudentsByStatus(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")

	var students []models.Student
	if err := database.DB.Where("status = ?", status).Find(&students).Error; err != nil {
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

	var updatedStudent models.Student
	if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedStudent.FirstName != "" {
		student.FirstName = updatedStudent.FirstName
	}
	if updatedStudent.LastName != "" {
		student.LastName = updatedStudent.LastName
	}
	if updatedStudent.Email != "" {
		student.Email = updatedStudent.Email
	}
	if updatedStudent.Status != "" {
		student.Status = updatedStudent.Status
	}

	if err := database.DB.Save(&student).Error; err != nil {
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

// SearchStudents searches students by name, email, or student ID
func SearchStudents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	var students []models.Student
	searchPattern := "%" + query + "%"
	if err := database.DB.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR student_id ILIKE ?", 
		searchPattern, searchPattern, searchPattern, searchPattern).Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
