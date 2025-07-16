package entities

import (
	"testing"

	"github.com/golang/011api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestUserInput_Validate(t *testing.T) {
	tests := []struct{
		name	string
		input	*entities.UserInput
		wantErr	bool
		errMsg	string
	}{
		{
			name: "valid user input",
			input: &entities.UserInput{
				Email: "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			input: &entities.UserInput{
				Email: "",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: true,
			errMsg: "email is required",
		},
		{
			name: "invalid email format",
			input: &entities.UserInput{
				Email: "invalid-email",
				Username: "testuser",
				Password: "password123",
			},
			wantErr: true,
			errMsg: "email must be valid",
		},
		{
			name: "empty username",
			input: &entities.UserInput{
				Email: "text@example.com",
				Username: "",
				Password: "password123",
			},
			wantErr: true,
			errMsg: "username is required",
		},
		{
			name: "username too short",
			input: &entities.UserInput{
				Email: "text@example.com",
				Username: "tu",
				Password: "password123",
			},
			wantErr: true,
			errMsg: "username must be less than 51 characters and more than 2 characters",
		},
		{
			name: "password too short",
			input: &entities.UserInput{
				Email: "text@example.com",
				Username: "testuser",
				Password: "pa",
			},
			wantErr: true,
			errMsg: "password must be more than 7 characters",
		},
		{
			name: "email and username with whitespace",
			input: &entities.UserInput{
				Email: "  TEST@EXAMPLE.COM  ",
				Username: "  testuser  ",
				Password: "password123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				if tt.input.Email == "  TEST@EXAMPLE.COM  " {
					assert.Equal(t, "test@example.com", tt.input.Email)
				}
				if tt.input.Username == "  testuser  " {
					assert.Equal(t, "testuser", tt.input.Username)
				}
			}
		})
	}
}