package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserData_JSONMarshaling(t *testing.T) {
	user := UserData{
		First_name:  "John",
		Second_name: "Doe",
		Email:       "john@example.com",
		UserRoleID:  1,
	}

	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var unmarshaledUser UserData
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	assert.NoError(t, err)
	assert.Equal(t, user.First_name, unmarshaledUser.First_name)
	assert.Equal(t, user.Email, unmarshaledUser.Email)
}

func TestUserData_AllFields(t *testing.T) {
	user := UserData{
		First_name:  "Jane",
		Second_name: "Smith",
		Email:       "jane@example.com",
		UserRoleID:  2,
	}

	assert.Equal(t, "Jane", user.First_name)
	assert.Equal(t, "Smith", user.Second_name)
	assert.Equal(t, "jane@example.com", user.Email)
	assert.Equal(t, uint(2), user.UserRoleID)
}
