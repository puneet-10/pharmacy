package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"pharmacy/models"
	"strconv"
)

// SignUpHandler handles the user sign-up process
func SignUpHandler(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Check that the required fields are provided
	if user.Name == "" || user.Phone == "" || user.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name, Email, and Password are required")
	}

	// Check if the email already exists
	token, existingUser, _ := models.AuthenticateUser(user.Phone, "")
	if existingUser != nil {
		return echo.NewHTTPError(http.StatusConflict, "Phone Number already in use")
	}

	// Create the user with the is_admin field
	token, newUser, err := models.CreateUser(user.Name, user.Email, user.Phone, user.Password, user.IsAdmin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	userResponse := models.UserResponse{
		PhoneNumber: newUser.Phone,
		Name:        newUser.Name,
		IsAdmin:     newUser.IsAdmin,
	}
	response := map[string]interface{}{
		"token": token,
		"user":  userResponse,
	}
	// Return the created user
	return c.JSON(http.StatusCreated, response)
}

// UpdateUserHandler handles the user update process
func UpdateUserHandler(c echo.Context) error {
	id := c.Param("id")
	var user models.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Check that required fields are provided
	if user.Name == "" || user.Phone == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name and Phone are required")
	}

	// Assume "updatedBy" is passed in the request headers (could be set via middleware based on logged-in user)
	updatedBy := c.Request().Header.Get("X-Updated-By")
	if updatedBy == "" {
		updatedBy = "system" // Default value if not provided
	}
	intId, _ := strconv.Atoi(id)
	// Update the user with the is_admin field
	updatedUser, err := models.UpdateUser(intId, user.Name, user.Phone, user.IsAdmin, updatedBy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Return the updated user
	return c.JSON(http.StatusOK, updatedUser)
}

// AuthenticateHandler authenticates a user based on email or phone number
func AuthenticateHandler(c echo.Context) error {
	var request struct {
		Identifier string `json:"identifier"` // Email or Phone
		Password   string `json:"password"`   // Password
	}

	// Bind request payload
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Ensure the identifier and password are provided
	if request.Identifier == "" || request.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Identifier and Password are required")
	}

	// Authenticate the user (either by email or phone)
	token, user, err := models.AuthenticateUser(request.Identifier, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// If the user is not found, return unauthorized error
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	userResponse := models.UserResponse{
		PhoneNumber: user.Phone,
		Name:        user.Name,
		IsAdmin:     user.IsAdmin,
	}
	response := map[string]interface{}{
		"token": token,
		"user":  userResponse,
	}
	// Return the authenticated user
	return c.JSON(http.StatusOK, response)
}
