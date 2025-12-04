package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStudentBankDetails_JSONMarshaling(t *testing.T) {
	id := uuid.New()
	studentID := uuid.New()
	bankID := uuid.New()

	details := StudentBankDetails{
		ID:                id,
		StudentID:         studentID,
		BankID:            bankID,
		AccountNumber:     "1234567890",
		AccountHolderName: "John Doe",
	}

	jsonData, err := json.Marshal(details)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), "John Doe")
	assert.Contains(t, string(jsonData), "1234567890")

	var unmarshaled StudentBankDetails
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", unmarshaled.AccountHolderName)
	assert.Equal(t, studentID, unmarshaled.StudentID)
}

func TestStudentBankDetails_JSONUnmarshaling(t *testing.T) {
	studentID := uuid.New()
	bankID := uuid.New()
	jsonStr := `{"student_id":"` + studentID.String() + `","bank_id":"` + bankID.String() + `","account_number":"1234567890","account_holder_name":"John Doe"}`

	var details StudentBankDetails
	err := json.Unmarshal([]byte(jsonStr), &details)
	assert.NoError(t, err)
	assert.Equal(t, studentID, details.StudentID)
	assert.Equal(t, bankID, details.BankID)
	assert.Equal(t, "1234567890", details.AccountNumber)
	assert.Equal(t, "John Doe", details.AccountHolderName)
}

func TestStudentBankDetails_ForeignKey(t *testing.T) {
	bankID := uuid.New()
	details := StudentBankDetails{
		BankID: bankID,
		Bank: Bank{
			ID:   bankID,
			Name: "Test Bank",
		},
	}

	assert.Equal(t, bankID, details.BankID)
	assert.Equal(t, "Test Bank", details.Bank.Name)
}

func TestStudentBankDetails_AllFields(t *testing.T) {
	id := uuid.New()
	studentID := uuid.New()
	bankID := uuid.New()

	details := StudentBankDetails{
		ID:                id,
		StudentID:         studentID,
		BankID:            bankID,
		AccountNumber:     "1234567890",
		AccountHolderName: "John Doe",
		Bank: Bank{
			ID:   bankID,
			Name: "Test Bank",
		},
	}

	assert.Equal(t, id, details.ID)
	assert.Equal(t, studentID, details.StudentID)
	assert.Equal(t, bankID, details.BankID)
	assert.Equal(t, "1234567890", details.AccountNumber)
	assert.Equal(t, "John Doe", details.AccountHolderName)
	assert.Equal(t, "Test Bank", details.Bank.Name)
}
