package models

import (
"encoding/json"
"testing"

"github.com/stretchr/testify/assert"
)

func TestUserRole_JSONMarshaling(t *testing.T) {
	role := UserRole{
		Name:        "Student",
		Description: "Student role",
	}

	jsonData, err := json.Marshal(role)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var unmarshaledRole UserRole
	err = json.Unmarshal(jsonData, &unmarshaledRole)
	assert.NoError(t, err)
	assert.Equal(t, role.Name, unmarshaledRole.Name)
}

func TestUserRole_AllFields(t *testing.T) {
	role := UserRole{
		Name:        "Admin",
		Description: "Administrator role",
	}

	assert.Equal(t, "Admin", role.Name)
	assert.Equal(t, "Administrator role", role.Description)
}
