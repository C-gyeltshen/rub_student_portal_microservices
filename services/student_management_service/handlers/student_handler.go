package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

// BulkCreateStudents creates multiple students from CSV import
func BulkCreateStudents(w http.ResponseWriter, r *http.Request) {
	var requestData []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(requestData) == 0 {
		http.Error(w, "No students provided", http.StatusBadRequest)
		return
	}

	// Process each student individually
	successCount := 0
	failedStudents := []map[string]interface{}{}
	successfulStudents := []models.Student{}
	
	for _, data := range requestData {
		student := models.Student{}
		
		// Extract basic fields
		if name, ok := data["last_name"].(string); ok {
			student.Name = name
		}
		if rubId, ok := data["rub_id_card_number"].(string); ok {
			student.RubIDCardNumber = rubId
		}
		if email, ok := data["email"].(string); ok {
			student.Email = email
		}
		if phone, ok := data["phone_number"].(string); ok {
			student.PhoneNumber = phone
		}
		if dob, ok := data["date_of_birth"].(string); ok {
			student.DateOfBirth = dob
		}
		
		// Validate required fields
		if student.RubIDCardNumber == "" || student.Name == "" || student.Email == "" {
			failedStudents = append(failedStudents, map[string]interface{}{
				"student": data,
				"error":   "Missing required fields (name, rub_id_card_number, email)",
			})
			continue
		}
		
		// Check if student already exists by email or rub_id_card_number
		var existing models.Student
		if err := database.DB.Where("email = ? OR rub_id_card_number = ?", student.Email, student.RubIDCardNumber).First(&existing).Error; err == nil {
			// Student already exists, skip it
			failedStudents = append(failedStudents, map[string]interface{}{
				"student": map[string]interface{}{
					"name":                student.Name,
					"rub_id_card_number":  student.RubIDCardNumber,
					"email":               student.Email,
				},
				"error": "Student already exists with this email or RUB ID",
			})
			continue
		}
		
		// Handle Program - lookup by name or use ID
		if programName, ok := data["program_name"].(string); ok && programName != "" {
			var program models.Program
			// Try to find existing program by name
			if err := database.DB.Where("name = ?", programName).First(&program).Error; err != nil {
				// Program not found, check if we have college info to create it
				if collegeName, ok := data["college_name"].(string); ok && collegeName != "" {
					var college models.College
					// Find or create college
					if err := database.DB.Where("name = ?", collegeName).FirstOrCreate(&college, models.College{Name: collegeName}).Error; err != nil {
						failedStudents = append(failedStudents, map[string]interface{}{
							"student": data,
							"error":   "Failed to find/create college: " + err.Error(),
						})
						continue
					}
					// Create program under this college
					program = models.Program{Name: programName, CollegeID: college.ID}
					if err := database.DB.Create(&program).Error; err != nil {
						failedStudents = append(failedStudents, map[string]interface{}{
							"student": data,
							"error":   "Failed to create program: " + err.Error(),
						})
						continue
					}
				} else {
					failedStudents = append(failedStudents, map[string]interface{}{
						"student": data,
						"error":   "Program '" + programName + "' not found and no college_name provided to create it",
					})
					continue
				}
			}
			student.ProgramID = program.ID
			student.CollegeID = program.CollegeID
		} else if programId, ok := data["program_id"].(float64); ok && programId > 0 {
			// Use program_id if provided
			student.ProgramID = uint(programId)
			// Get college from program
			var program models.Program
			if err := database.DB.First(&program, programId).Error; err == nil {
				student.CollegeID = program.CollegeID
			}
		}
		
		// Handle College separately if provided and not set by program
		if student.CollegeID == 0 {
			if collegeName, ok := data["college_name"].(string); ok && collegeName != "" {
				var college models.College
				if err := database.DB.Where("name = ?", collegeName).FirstOrCreate(&college, models.College{Name: collegeName}).Error; err != nil {
					failedStudents = append(failedStudents, map[string]interface{}{
						"student": data,
						"error":   "Failed to find/create college: " + err.Error(),
					})
					continue
				}
				student.CollegeID = college.ID
			} else if collegeId, ok := data["college_id"].(float64); ok && collegeId > 0 {
				student.CollegeID = uint(collegeId)
			}
		}
		
		// Handle user_id
		if userId, ok := data["user_id"].(float64); ok && userId > 0 {
			student.UserID = uint(userId)
		}
		
		// Check if student already exists by email or rub_id_card_number
		var existingStudent models.Student
		err := database.DB.Where("email = ? OR rub_id_card_number = ?", student.Email, student.RubIDCardNumber).First(&existingStudent).Error
		
		if err == nil {
			// Student already exists, skip it
			failedStudents = append(failedStudents, map[string]interface{}{
				"student": map[string]interface{}{
					"name":                student.Name,
					"rub_id_card_number":  student.RubIDCardNumber,
					"email":               student.Email,
				},
				"error": "Student already exists with this email or ID card number (skipped)",
			})
			continue
		}
		
		// Create student (only if not exists)
		if err := database.DB.Create(&student).Error; err != nil {
			failedStudents = append(failedStudents, map[string]interface{}{
				"student": map[string]interface{}{
					"name":                student.Name,
					"rub_id_card_number":  student.RubIDCardNumber,
					"email":               student.Email,
				},
				"error": err.Error(),
			})
		} else {
			successCount++
			// Reload with relationships
			database.DB.Preload("Program").Preload("College").First(&student, student.ID)
			successfulStudents = append(successfulStudents, student)
		}
	}

	response := map[string]interface{}{
		"success":             successCount,
		"failed":              len(failedStudents),
		"total":               len(requestData),
		"failed_records":      failedStudents,
		"successful_records":  successfulStudents,
	}

	w.Header().Set("Content-Type", "application/json")
	if len(failedStudents) > 0 && successCount > 0 {
		w.WriteHeader(http.StatusOK) // Partial success
	} else if successCount == 0 {
		w.WriteHeader(http.StatusBadRequest) // All failed
	} else {
		w.WriteHeader(http.StatusCreated) // All succeeded
	}
	json.NewEncoder(w).Encode(response)
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

// BulkCreateStudents creates multiple students at once
func BulkCreateStudents(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	if err := json.NewDecoder(r.Body).Decode(&students); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to parse JSON: " + err.Error(),
			"error":   err.Error(),
		})
		return
	}

	if len(students) == 0 {
		http.Error(w, "No students provided", http.StatusBadRequest)
		return
	}

	// Validate all students with detailed error reporting
	var validationErrors []map[string]interface{}
	for i, student := range students {
		missingFields := []string{}
		
		if student.StudentID == "" {
			missingFields = append(missingFields, "student_id")
		}
		if student.FirstName == "" {
			missingFields = append(missingFields, "first_name")
		}
		if student.LastName == "" {
			missingFields = append(missingFields, "last_name")
		}
		if student.Email == "" {
			missingFields = append(missingFields, "email")
		}
		
		if len(missingFields) > 0 {
			validationErrors = append(validationErrors, map[string]interface{}{
				"index":          i,
				"missing_fields": missingFields,
				"received_data": map[string]interface{}{
					"student_id": student.StudentID,
					"first_name": student.FirstName,
					"last_name":  student.LastName,
					"email":      student.Email,
					"cid":        student.CID,
					"gender":     student.Gender,
				},
			})
		}
	}
	
	// If there are validation errors, return them all
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":           false,
			"message":           fmt.Sprintf("Validation failed for %d student(s)", len(validationErrors)),
			"validation_errors": validationErrors,
			"total_students":    len(students),
		})
		return
	}

	// Create all students in a transaction
	tx := database.DB.Begin()
	
	successCount := 0
	var errors []string
	
	for _, student := range students {
		if err := tx.Create(&student).Error; err != nil {
			errors = append(errors, "Failed to create student "+student.StudentID+": "+err.Error())
		} else {
			successCount++
		}
	}

	if len(errors) > 0 && successCount == 0 {
		tx.Rollback()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to create any students",
			"errors":  errors,
		})
		return
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"created_count": successCount,
		"total_count":   len(students),
		"errors":        errors,
		"message":       "Bulk creation completed",
	})
}

// BulkCreateStudentsFromCSV handles CSV file upload and creates students
func BulkCreateStudentsFromCSV(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to parse form: " + err.Error(),
		})
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to get file from form: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	
	// Read header row
	headers, err := reader.Read()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to read CSV headers: " + err.Error(),
		})
		return
	}

	// Create a map for header indices
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	var students []models.Student
	var parseErrors []map[string]interface{}
	rowNum := 1 // Start at 1 (header is row 0)

	// Read all rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			parseErrors = append(parseErrors, map[string]interface{}{
				"row":     rowNum,
				"message": "Failed to read row: " + err.Error(),
			})
			rowNum++
			continue
		}

		// Helper function to get value from record
		getValue := func(fieldName string) string {
			if idx, ok := headerMap[fieldName]; ok && idx < len(record) {
				return record[idx]
			}
			return ""
		}

		// Parse student from CSV row
		student := models.Student{
			StudentID:           getValue("student_id"),
			FirstName:           getValue("first_name"),
			LastName:            getValue("last_name"),
			Email:               getValue("email"),
			PhoneNumber:         getValue("phone_number"),
			CID:                 getValue("cid"),
			DateOfBirth:         getValue("date_of_birth"),
			Gender:              getValue("gender"),
			PermanentAddress:    getValue("permanent_address"),
			CurrentAddress:      getValue("current_address"),
			EnrollmentDate:      getValue("enrollment_date"),
			GraduationDate:      getValue("graduation_date"),
			Status:              getValue("status"),
			AcademicStanding:    getValue("academic_standing"),
			GuardianName:        getValue("guardian_name"),
			GuardianPhoneNumber: getValue("guardian_phone_number"),
			GuardianRelation:    getValue("guardian_relation"),
		}

		// Parse numeric fields with error handling
		if programID := getValue("program_id"); programID != "" {
			if id, err := strconv.ParseUint(programID, 10, 32); err == nil {
				student.ProgramID = uint(id)
			}
		}

		if collegeID := getValue("college_id"); collegeID != "" {
			if id, err := strconv.ParseUint(collegeID, 10, 32); err == nil {
				student.CollegeID = uint(id)
			}
		}

		if yearOfStudy := getValue("year_of_study"); yearOfStudy != "" {
			if year, err := strconv.Atoi(yearOfStudy); err == nil {
				student.YearOfStudy = year
			}
		}

		if semester := getValue("semester"); semester != "" {
			if sem, err := strconv.Atoi(semester); err == nil {
				student.Semester = sem
			}
		}

		if gpa := getValue("gpa"); gpa != "" {
			if gpaVal, err := strconv.ParseFloat(gpa, 64); err == nil {
				student.GPA = gpaVal
			}
		}

		students = append(students, student)
		rowNum++
	}

	// Validate all students
	var validationErrors []map[string]interface{}
	for i, student := range students {
		missingFields := []string{}
		
		if student.StudentID == "" {
			missingFields = append(missingFields, "student_id")
		}
		if student.FirstName == "" {
			missingFields = append(missingFields, "first_name")
		}
		if student.LastName == "" {
			missingFields = append(missingFields, "last_name")
		}
		if student.Email == "" {
			missingFields = append(missingFields, "email")
		}
		
		if len(missingFields) > 0 {
			validationErrors = append(validationErrors, map[string]interface{}{
				"row":            i + 2, // +2 because row 0 is header, row indices start at 0
				"missing_fields": missingFields,
				"received_data": map[string]interface{}{
					"student_id": student.StudentID,
					"first_name": student.FirstName,
					"last_name":  student.LastName,
					"email":      student.Email,
				},
			})
		}
	}
	
	// If there are validation errors, return them
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":           false,
			"message":           fmt.Sprintf("Validation failed for %d student(s)", len(validationErrors)),
			"validation_errors": validationErrors,
			"parse_errors":      parseErrors,
			"total_rows":        len(students),
		})
		return
	}

	// Create all students in a transaction
	tx := database.DB.Begin()
	
	successCount := 0
	var errors []string
	
	for i, student := range students {
		if err := tx.Create(&student).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Row %d (%s): %s", i+2, student.StudentID, err.Error()))
		} else {
			successCount++
		}
	}

	if len(errors) > 0 && successCount == 0 {
		tx.Rollback()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to create any students",
			"errors":  errors,
		})
		return
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"created_count": successCount,
		"total_count":   len(students),
		"errors":        errors,
		"parse_errors":  parseErrors,
		"message":       "CSV bulk creation completed",
	})
}

