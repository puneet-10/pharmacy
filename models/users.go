package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

// User struct represents the user model in the database
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone" gorm:"unique"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

type UserResponse struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	IsAdmin     bool   `json:"isAdmin"`
}

var db *gorm.DB
var jwtSecret = []byte("RaghavSureka")

// SetDB initializes the database connection for GORM
func SetDB(d *gorm.DB) {
	db = d
}

// GenerateJWT creates a JWT token for the given user
func GenerateJWT(user *User) (string, error) {
	claims := jwt.MapClaims{
		"userId":  user.ID, // <-- Include user ID
		"name":    user.Name,
		"phone":   user.Phone,
		"isAdmin": user.IsAdmin,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// CreateUser creates a new user in the database
func CreateUser(name, email, phone, password string, isAdmin bool) (string, *User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	user := &User{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
		IsAdmin:  isAdmin,
	}

	if err := db.Create(user).Error; err != nil {
		return "", nil, err
	}

	token, err := GenerateJWT(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// AuthenticateUser authenticates a user by phone and password
func AuthenticateUser(identifier, password string) (string, *User, error) {
	var user User

	if err := db.Where("phone = ?", identifier).First(&user).Error; err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, err // Invalid password
	}

	token, err := GenerateJWT(&user)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

// UpdateUser updates a user's data (name, phone, is_admin, updated_at)
func UpdateUser(id int, name, phone string, isAdmin bool, updatedBy string) (*User, error) {
	var user User

	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	user.Name = name
	user.Phone = phone
	user.IsAdmin = isAdmin
	user.UpdatedAt = time.Now()
	user.UpdatedBy = updatedBy

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Helper function to check if a string is an email
func isEmail(str string) bool {
	return strings.Contains(str, "@") && strings.Contains(str, ".")
}

// AuthMiddleware decodes JWT token and sets user context
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.ErrUnauthorized
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return echo.ErrUnauthorized
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_name", claims["name"])
		c.Set("user_phoneNumber", claims["phone"])
		c.Set("user_isAdmin", claims["is_admin"])

		return next(c)
	}
}