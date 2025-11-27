package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"student_management_service/database"
	client "student_management_service/grpc/client"
	"student_management_service/models"
	"time"

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

	// Connect to finance service to calculate stipend with deductions
	financeClient, err := client.NewFinanceClient()
	if err != nil {
		log.Printf("Warning: Could not connect to finance service: %v", err)
		// Continue without finance service integration
	} else {
		defer financeClient.Close()

		// Calculate stipend with deductions using finance service
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		calcResult, err := financeClient.CalculateStipendWithDeductions(
			ctx,
			fmt.Sprintf("%d", student.ID),
			student.FinancingType,
			allocation.Amount,
		)
		if err != nil {
			log.Printf("Warning: Stipend calculation failed: %v", err)
		} else {
			// Update allocation with calculated net amount
			allocation.Amount = calcResult.NetStipendAmount
			log.Printf("Stipend calculated: Base=%.2f, Deductions=%.2f, Net=%.2f",
				calcResult.BaseStipendAmount,
				calcResult.TotalDeductions,
				calcResult.NetStipendAmount)
		}
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

	// Get student to determine stipend type
	var student models.Student
	if err := database.DB.First(&student, history.StudentID).Error; err != nil {
		log.Printf("Warning: Could not find student %d: %v", history.StudentID, err)
	} else {
		// Also create stipend record in finance service
		financeClient, err := client.NewFinanceClient()
		if err != nil {
			log.Printf("Warning: Could not connect to finance service: %v", err)
		} else {
			defer financeClient.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Determine stipend type from student financing type
			stipendType := student.FinancingType
			if stipendType == "" {
				stipendType = "scholarship"
			}

			// Create stipend in finance service
			_, err = financeClient.CreateStipend(
				ctx,
				fmt.Sprintf("%d", history.StudentID),
				stipendType,
				history.Amount,
				history.PaymentMethod,
				history.BankReference, // Using BankReference as JournalNumber
				history.Remarks,
			)
			if err != nil {
				log.Printf("Warning: Failed to create stipend in finance service: %v", err)
			} else {
				log.Printf("Successfully created stipend in finance service for student %d", history.StudentID)
			}
		}
	}

	if err := database.DB.Create(&history).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(history)
}

// ==================== Finance Service Integration Endpoints ====================

// CalculateStipendWithDeductions calculates stipend with deductions via finance service
func CalculateStipendWithDeductions(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StudentID   uint    `json:"student_id"`
		StipendType string  `json:"stipend_type"`
		Amount      float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate student exists
	var student models.Student
	if err := database.DB.First(&student, req.StudentID).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Use student's financing type if not provided
	stipendType := req.StipendType
	if stipendType == "" {
		stipendType = student.FinancingType
	}

	// Connect to finance service
	financeClient, err := client.NewFinanceClient()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to finance service: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer financeClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate stipend with deductions
	result, err := financeClient.CalculateStipendWithDeductions(
		ctx,
		fmt.Sprintf("%d", req.StudentID),
		stipendType,
		req.Amount,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Calculation failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetStudentFinanceStipends retrieves stipends from finance service for a student
func GetStudentFinanceStipends(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")

	// Validate student exists
	id, err := strconv.Atoi(studentID)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var student models.Student
	if err := database.DB.First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Connect to finance service
	financeClient, err := client.NewFinanceClient()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to finance service: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer financeClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get student stipends from finance service
	result, err := financeClient.GetStudentStipends(ctx, studentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get stipends: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
