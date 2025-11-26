package models

import (
"encoding/json"
"testing"
"github.com/stretchr/testify/assert"
"gorm.io/gorm"
)

func TestBank_JSONMarshaling(t *testing.T) {
	bank := Bank{
		Model: gorm.Model{ID: 1},
		Name:  "Test Bank",
	}

	jsonData, err := json.Marshal(bank)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), "Test Bank")

	var unmarshaledBank Bank
	err = json.Unmarshal(jsonData, &unmarshaledBank)
	assert.NoError(t, err)
	assert.Equal(t, "Test Bank", unmarshaledBank.Name)
}

func TestBank_Relationships(t *testing.T) {
	bank := Bank{
		Model:              gorm.Model{ID: 1},
		Name:               "Test Bank",
		StudentBankDetails: []StudentBankDetails{},
	}

	assert.NotNil(t, bank.StudentBankDetails)
	assert.Len(t, bank.StudentBankDetails, 0)
}
