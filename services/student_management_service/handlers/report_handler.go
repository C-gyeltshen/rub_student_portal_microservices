package handlers

import (
	"encoding/json"
	"net/http"
	"student_management_service/database"
	"student_management_service/models"
)

// StudentSummary represents aggregated student data
type StudentSummary struct {
	TotalStudents      int                `json:"total_students"`
	ActiveStudents     int                `json:"active_students"`
	InactiveStudents   int                `json:"inactive_students"`
	GraduatedStudents  int                `json:"graduated_students"`
	ByCollege          map[string]int     `json:"by_college"`
	ByProgram          map[string]int     `json:"by_program"`
	ByStatus           map[string]int     `json:"by_status"`
	StipendEligible    int                `json:"stipend_eligible"`
	StipendIneligible  int                `json:"stipend_ineligible"`
}

// StipendStatistics represents stipend-related statistics
type StipendStatistics struct {
	TotalAllocations     int     `json:"total_allocations"`
	TotalAmount          float64 `json:"total_amount"`
	PendingAllocations   int     `json:"pending_allocations"`
	ApprovedAllocations  int     `json:"approved_allocations"`
	DisbursedAmount      float64 `json:"disbursed_amount"`
	ByProgram            map[string]float64 `json:"by_program"`
	ByCollege            map[string]float64 `json:"by_college"`
}

// GenerateStudentSummary generates a summary report of all students
func GenerateStudentSummary(w http.ResponseWriter, r *http.Request) {
	summary := StudentSummary{
		ByCollege: make(map[string]int),
		ByProgram: make(map[string]int),
		ByStatus:  make(map[string]int),
	}

	// Get total students
	var students []models.Student
	if err := database.DB.Preload("College").Preload("Program").Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	summary.TotalStudents = len(students)

	// Calculate statistics
	for _, student := range students {
		// Count by status
		switch student.Status {
		case "active":
			summary.ActiveStudents++
		case "inactive":
			summary.InactiveStudents++
		case "graduated":
			summary.GraduatedStudents++
		}
		summary.ByStatus[student.Status]++

		// Count by college
		if student.College.ID != 0 {
			summary.ByCollege[student.College.Name]++
		}

		// Count by program
		if student.Program.ID != 0 {
			summary.ByProgram[student.Program.Name]++
		}

		// Check stipend eligibility
		eligibility := calculateEligibility(student)
		if eligibility.IsEligible {
			summary.StipendEligible++
		} else {
			summary.StipendIneligible++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GenerateStipendStatistics generates stipend allocation statistics
func GenerateStipendStatistics(w http.ResponseWriter, r *http.Request) {
	stats := StipendStatistics{
		ByProgram: make(map[string]float64),
		ByCollege: make(map[string]float64),
	}

	// Get all allocations
	var allocations []models.StipendAllocation
	if err := database.DB.Find(&allocations).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats.TotalAllocations = len(allocations)

	for _, allocation := range allocations {
		stats.TotalAmount += allocation.Amount

		switch allocation.Status {
		case "pending":
			stats.PendingAllocations++
		case "approved", "disbursed":
			stats.ApprovedAllocations++
			if allocation.Status == "disbursed" {
				stats.DisbursedAmount += allocation.Amount
			}
		}

		// Get student with program and college info
		var student models.Student
		if err := database.DB.Preload("Program").Preload("College").First(&student, allocation.StudentID).Error; err == nil {
			// Aggregate by program
			if student.Program.ID != 0 {
				stats.ByProgram[student.Program.Name] += allocation.Amount
			}

			// Aggregate by college
			if student.College.ID != 0 {
				stats.ByCollege[student.College.Name] += allocation.Amount
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetStudentsByCollegeReport gets detailed student list by college
func GetStudentsByCollegeReport(w http.ResponseWriter, r *http.Request) {
	collegeID := r.URL.Query().Get("college_id")
	if collegeID == "" {
		http.Error(w, "college_id parameter required", http.StatusBadRequest)
		return
	}

	var students []models.Student
	if err := database.DB.Where("college_id = ?", collegeID).
		Preload("Program").Preload("College").
		Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetStudentsByProgramReport gets detailed student list by program
func GetStudentsByProgramReport(w http.ResponseWriter, r *http.Request) {
	programID := r.URL.Query().Get("program_id")
	if programID == "" {
		http.Error(w, "program_id parameter required", http.StatusBadRequest)
		return
	}

	var students []models.Student
	if err := database.DB.Where("program_id = ?", programID).
		Preload("Program").Preload("College").
		Find(&students).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
