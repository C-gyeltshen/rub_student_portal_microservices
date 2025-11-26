package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"student_management_service/database"
	"student_management_service/models"

	"github.com/go-chi/chi/v5"
)

// ==================== Stipend Eligibility ====================

// CheckStipendEligibility checks if a student is eligible for stipend
func CheckStipendEligibility(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")
	id, err := strconv.Atoi(studentID)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var student models.Student
	if err := database.DB.Preload("Program").First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	eligibility := calculateEligibility(student)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eligibility)
}

// calculateEligibility determines student eligibility for stipend
func calculateEligibility(student models.Student) models.StipendEligibility {
	eligibility := models.StipendEligibility{
		StudentID: student.ID,
		Reasons:   []string{},
	}

	// Check if student is active
	if student.Status != "active" {
		eligibility.IsEligible = false
		eligibility.Reasons = append(eligibility.Reasons, "Student is not active")
		return eligibility
	}

	// Preload college to check stipend policy
	database.DB.Preload("College").First(&student, student.ID)

	// Check financing type and college policy
	if student.FinancingType == "self-financed" {
		// Check if college allows self-financed students to get stipend
		if student.College.ID != 0 && !student.College.AllowSelfFinancedStipend {
			eligibility.IsEligible = false
			eligibility.Reasons = append(eligibility.Reasons, "College does not allow self-financed students to receive stipend")
			return eligibility
		}
	}
	// Scholarship students are always eligible for stipend (if other criteria met)

	// Check if program has stipend
	if student.Program.ID != 0 && !student.Program.HasStipend {
		eligibility.IsEligible = false
		eligibility.Reasons = append(eligibility.Reasons, "Program does not offer stipend")
		return eligibility
	}

	// Student is eligible
	eligibility.IsEligible = true
	eligibility.ExpectedAmount = student.Program.StipendAmount
	eligibility.FinancingType = student.FinancingType
	eligibility.Reasons = append(eligibility.Reasons, "Student meets all eligibility criteria")

	return eligibility
}

// ==================== Stipend Allocation ====================

// CreateStipendAllocation creates a new stipend allocation
func CreateStipendAllocation(w http.ResponseWriter, r *http.Request) {
	var allocation models.StipendAllocation
	if err := json.NewDecoder(r.Body).Decode(&allocation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate student exists and is eligible
	var student models.Student
	if err := database.DB.Preload("Program").First(&student, allocation.StudentID).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	eligibility := calculateEligibility(student)
	if !eligibility.IsEligible {
		http.Error(w, "Student is not eligible for stipend", http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&allocation).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(allocation)
}

// GetStipendAllocations retrieves all stipend allocations
func GetStipendAllocations(w http.ResponseWriter, r *http.Request) {
	var allocations []models.StipendAllocation
	
	query := database.DB.Preload("Student")
	
	// Filter by status if provided
	status := r.URL.Query().Get("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	// Filter by student if provided
	studentID := r.URL.Query().Get("student_id")
	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	
	if err := query.Find(&allocations).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allocations)
}

// GetStipendAllocationById retrieves a specific allocation
func GetStipendAllocationById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var allocation models.StipendAllocation
	if err := database.DB.Preload("Student").First(&allocation, id).Error; err != nil {
		http.Error(w, "Allocation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allocation)
}

// UpdateStipendAllocation updates an allocation (e.g., approve/reject)
func UpdateStipendAllocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var allocation models.StipendAllocation
	if err := database.DB.First(&allocation, id).Error; err != nil {
		http.Error(w, "Allocation not found", http.StatusNotFound)
		return
	}

	var updates models.StipendAllocation
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update allowed fields
	if updates.Status != "" {
		allocation.Status = updates.Status
	}
	if updates.ApprovedBy != 0 {
		allocation.ApprovedBy = updates.ApprovedBy
	}
	if updates.ApprovalDate != "" {
		allocation.ApprovalDate = updates.ApprovalDate
	}
	if updates.Remarks != "" {
		allocation.Remarks = updates.Remarks
	}

	if err := database.DB.Save(&allocation).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allocation)
}

// ==================== Stipend History ====================

// GetStipendHistory retrieves stipend payment history
func GetStipendHistory(w http.ResponseWriter, r *http.Request) {
	var history []models.StipendHistory
	
	query := database.DB.Preload("Student")
	
	// Filter by student if provided
	studentID := r.URL.Query().Get("student_id")
	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	
	// Filter by status if provided
	status := r.URL.Query().Get("status")
	if status != "" {
		query = query.Where("transaction_status = ?", status)
	}
	
	if err := query.Order("payment_date DESC").Find(&history).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// GetStudentStipendHistory retrieves stipend history for a specific student
func GetStudentStipendHistory(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")

	var history []models.StipendHistory
	if err := database.DB.Where("student_id = ?", studentID).
		Order("payment_date DESC").
		Find(&history).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// CreateStipendHistory records a stipend payment
func CreateStipendHistory(w http.ResponseWriter, r *http.Request) {
	var history models.StipendHistory
	if err := json.NewDecoder(r.Body).Decode(&history); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&history).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(history)
}
