package models

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestStudentBankDetails_JSONMarshaling(t *testing.T) {
	details := StudentBankDetails{
		Model:             gorm.Model{ID: 1},
		StudentID:         123,
		BankID:            1,
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
	assert.Equal(t, 123, unmarshaled.StudentID)
}

func TestStudentBankDetails_JSONUnmarshaling(t *testing.T) {
	jsonStr := `{"student_id":123,"bank_id":1,"account_number":"1234567890","account_holder_name":"John Doe"}`

	var details StudentBankDetails
	err := json.Unmarshal([]byte(jsonStr), &details)
	assert.NoError(t, err)
	assert.Equal(t, 123, details.StudentID)
	assert.Equal(t, uint(1), details.BankID)
	assert.Equal(t, "1234567890", details.AccountNumber)
	assert.Equal(t, "John Doe", details.AccountHolderName)
}

func TestStudentBankDetails_ForeignKey(t *testing.T) {
	details := StudentBankDetails{
		BankID: 1,
		Bank: Bank{
			Model: gorm.Model{ID: 1},
			Name:  "Test Bank",
		},
	}

	assert.Equal(t, uint(1), details.BankID)
	assert.Equal(t, "Test Bank", details.Bank.Name)
}

func TestStudentBankDetails_AllFields(t *testing.T) {
	details := StudentBankDetails{
		Model:             gorm.Model{ID: 1},
		StudentID:         123,
		BankID:            1,
		AccountNumber:     "1234567890",
		AccountHolderName: "John Doe",
		Bank: Bank{
			Model: gorm.Model{ID: 1},
			Name:  "Test Bank",
		},
	}

	assert.Equal(t, uint(1), details.ID)
	assert.Equal(t, 123, details.StudentID)
	assert.Equal(t, uint(1), details.BankID)
	assert.Equal(t, "1234567890", details.AccountNumber)
	assert.Equal(t, "John Doe", details.AccountHolderName)
	assert.Equal(t, "Test Bank", details.Bank.Name)
}
