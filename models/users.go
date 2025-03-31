package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

// User struct represents the user model in the database
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Phone     string    `json:"phone" gorm:"unique"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

var db *gorm.DB

// SetDB initializes the database connection for GORM
func SetDB(d *gorm.DB) {
	db = d
}

// CreateUser creates a new user in the database
func CreateUser(name, email, phone, password string, isAdmin bool) (*User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create a new user instance
	user := &User{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
		IsAdmin:  isAdmin,
	}

	// Insert the user into the database
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	// Return the created user
	return user, nil
}

// AuthenticateUser authenticates a user by email or phone and password
func AuthenticateUser(identifier, password string) (*User, error) {
	var user User

	// Check if the identifier is email or phone
	if isEmail(identifier) {
		// If it's an email, search by email
		if err := db.Where("email = ?", identifier).First(&user).Error; err != nil {
			return nil, err
		}
	} else {
		// Otherwise, treat it as a phone number and search by phone
		if err := db.Where("phone = ?", identifier).First(&user).Error; err != nil {
			return nil, err
		}
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err // Invalid password
	}

	return &user, nil
}

// UpdateUser updates a user's data (name, phone, is_admin, updated_at)
func UpdateUser(id int, name, phone string, isAdmin bool, updatedBy string) (*User, error) {
	var user User

	// Find the user by ID
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Update the user's fields
	user.Name = name
	user.Phone = phone
	user.IsAdmin = isAdmin
	user.UpdatedAt = time.Now()
	user.UpdatedBy = updatedBy

	// Save the updated user to the database
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Helper function to check if a string is an email
func isEmail(str string) bool {
	// Simple regex check for email-like string
	return strings.Contains(str, "@") && strings.Contains(str, ".")
}
