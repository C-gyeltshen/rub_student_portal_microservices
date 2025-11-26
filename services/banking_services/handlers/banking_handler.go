package handlers

import (
	"banking_services/database"
	"banking_services/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ==================== Bank Handlers ====================

// GetBanks retrieves all banks
func GetBanks(w http.ResponseWriter, r *http.Request) {
	var banks []models.Bank
	if err := database.DB.Find(&banks).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banks)
}

// CreateBank creates a new bank
func CreateBank(w http.ResponseWriter, r *http.Request) {
	var bank models.Bank
	if err := json.NewDecoder(r.Body).Decode(&bank); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&bank).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bank)
}

// GetBankById retrieves a bank by ID
func GetBankById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var bank models.Bank
	if err := database.DB.First(&bank, id).Error; err != nil {
		http.Error(w, "Bank not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bank)
}

// UpdateBank updates a bank by ID
func UpdateBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var bank models.Bank
	if err := database.DB.First(&bank, id).Error; err != nil {
		http.Error(w, "Bank not found", http.StatusNotFound)
		return
	}

	var updatedBank models.Bank
	if err := json.NewDecoder(r.Body).Decode(&updatedBank); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update only the fields that should be updated
	bank.Name = updatedBank.Name

	if err := database.DB.Save(&bank).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bank)
}

// DeleteBank deletes a bank by ID
func DeleteBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var bank models.Bank
	if err := database.DB.First(&bank, id).Error; err != nil {
		http.Error(w, "Bank not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Delete(&bank).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bank deleted successfully"})
}

// ==================== Student Bank Details Handlers ====================

// GetStudentBankDetails retrieves all student bank details
func GetStudentBankDetails(w http.ResponseWriter, r *http.Request) {
	var studentBankDetails []models.StudentBankDetails
	// Preload the Bank relationship to include bank information
	if err := database.DB.Preload("Bank").Find(&studentBankDetails).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentBankDetails)
}

// CreateStudentBankDetails creates new student bank details
func CreateStudentBankDetails(w http.ResponseWriter, r *http.Request) {
	var studentBankDetails models.StudentBankDetails
	if err := json.NewDecoder(r.Body).Decode(&studentBankDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify that the bank exists
	var bank models.Bank
	if err := database.DB.Where("id = ?", studentBankDetails.BankID).First(&bank).Error; err != nil {
		http.Error(w, "Bank not found", http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&studentBankDetails).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with Bank relationship
	database.DB.Preload("Bank").First(&studentBankDetails, studentBankDetails.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(studentBankDetails)
}

// GetStudentBankDetailsById retrieves student bank details by ID
func GetStudentBankDetailsById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var studentBankDetails models.StudentBankDetails
	if err := database.DB.Preload("Bank").First(&studentBankDetails, id).Error; err != nil {
		http.Error(w, "Student bank details not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentBankDetails)
}

// GetStudentBankDetailsByStudentId retrieves student bank details by student ID
func GetStudentBankDetailsByStudentId(w http.ResponseWriter, r *http.Request) {
	studentId := chi.URLParam(r, "studentId")

	var studentBankDetails []models.StudentBankDetails
	if err := database.DB.Preload("Bank").Where("student_id = ?", studentId).Find(&studentBankDetails).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(studentBankDetails) == 0 {
		http.Error(w, "No bank details found for this student", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentBankDetails)
}

// UpdateStudentBankDetails updates student bank details by ID
func UpdateStudentBankDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var studentBankDetails models.StudentBankDetails
	if err := database.DB.First(&studentBankDetails, id).Error; err != nil {
		http.Error(w, "Student bank details not found", http.StatusNotFound)
		return
	}

	var updatedDetails models.StudentBankDetails
	if err := json.NewDecoder(r.Body).Decode(&updatedDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify that the bank exists if BankID is being updated
	if updatedDetails.BankID != "" && updatedDetails.BankID != studentBankDetails.BankID {
		var bank models.Bank
		if err := database.DB.Where("id = ?", updatedDetails.BankID).First(&bank).Error; err != nil {
			http.Error(w, "Bank not found", http.StatusBadRequest)
			return
		}
	}

	// Update fields
	if updatedDetails.StudentID != "" {
		studentBankDetails.StudentID = updatedDetails.StudentID
	}
	if updatedDetails.BankID != "" {
		studentBankDetails.BankID = updatedDetails.BankID
	}
	if updatedDetails.AccountNumber != "" {
		studentBankDetails.AccountNumber = updatedDetails.AccountNumber
	}
	if updatedDetails.AccountHolderName != "" {
		studentBankDetails.AccountHolderName = updatedDetails.AccountHolderName
	}

	if err := database.DB.Save(&studentBankDetails).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with Bank relationship
	database.DB.Preload("Bank").First(&studentBankDetails, studentBankDetails.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studentBankDetails)
}

// DeleteStudentBankDetails deletes student bank details by ID
func DeleteStudentBankDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var studentBankDetails models.StudentBankDetails
	if err := database.DB.First(&studentBankDetails, id).Error; err != nil {
		http.Error(w, "Student bank details not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Delete(&studentBankDetails).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Student bank details deleted successfully"})
}
