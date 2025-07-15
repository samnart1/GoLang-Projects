package entities

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID				int			`json:"id"`
	Email			string		`json:"email"`
	Username		string		`json:"username"`
	PasswordHash	string		`json:"-"`
	FirstName		*string		`json:"first_name,omitempty"`
	LastName		*string		`json:"last_name,omitempty"`
	IsActive		bool		`json:"is_active"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}

type UserInput struct {
	Email		string	`json:"email" validate:"required,email"`
	Username	string	`json:"username" validate:"required,min=3,max=50"`
	Password	string	`json:"password" validate:"required,min=8"`
	FirstName	*string	`json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName	*string	`json:"last_name,omitempty" validate:"omitempty,max=100"`
}

func (u *UserInput) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	if u.Username == "" {
		return fmt.Errorf("username is required")
	}

	if u.Password == "" {
		return fmt.Errorf("password is required")
	}

	if !isValidEmail(u.Email) {
		return fmt.Errorf("email must be valid")
	}

	if len(u.Username) < 3 || len(u.Username) > 50 {
		return fmt.Errorf("username must be less than 51 characters and more than 2 characters")
	}

	if len(u.Password) < 8 {
		return fmt.Errorf("password must be more than 7 characters")
	}

	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Username = strings.TrimSpace(u.Username)

	return nil
}

func (u *UserInput) FullName() string {
	parts := []string{}

	if u.FirstName != nil && *u.FirstName != "" {
		parts = append(parts, *u.FirstName)
	}

	if u.LastName != nil && *u.LastName != "" {
		parts = append(parts, *u.LastName)
	}

	return strings.Join(parts, " ")
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}