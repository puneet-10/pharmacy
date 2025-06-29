package handlers

import (
	"errors"
	"net/http"
	"pharmacy/models"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct{}

// NewOrderHandler creates a new instance of the handler
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userID, _, err := GetUserFromHeader(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}

	var req models.OrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// ✅ Set userID from token
	req.UserID = userID

	order, err := models.CreateOrderWithItems(req, "api_user") // or fetch updatedBy from token
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create order"})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	// Check if user is admin to include user details
	_, isAdmin, err := GetUserFromHeader(c)
	includeUserDetails := false
	if err == nil && isAdmin {
		includeUserDetails = true
	}

	order, err := models.GetOrderByIDWithUserDetails(uint(id), includeUserDetails)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Order not found"})
	}
	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	userID, isAdmin, err := GetUserFromHeader(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}
	orders, err := models.GetAllOrders(userID, isAdmin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch orders"})
	}
	if orders == nil {
		return c.JSON(http.StatusOK, []models.OrderRequest{})
	}
	return c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	var req models.OrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	order, err := models.UpdateOrder(uint(id), req, "admin_user")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update order"})
	}
	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	// Validate status values
	validStatuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if req.Status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid status. Valid statuses are: pending, processing, shipped, delivered, cancelled"})
	}

	userID, _, err := GetUserFromHeader(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}

	updatedBy := "user_" + strconv.Itoa(int(userID))
	order, err := models.UpdateOrderStatus(uint(id), req.Status, updatedBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update order status"})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}
	if err := models.DeleteOrder(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete order"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Order deleted successfully"})
}

func GetUserFromToken(c echo.Context) (uint, bool, error) {
	userToken, ok := c.Get("token").(*jwt.Token)
	if !ok {
		return 0, false, errors.New("invalid token format")
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false, errors.New("invalid JWT claims")
	}

	userIDFloat, ok := claims["userId"].(float64)
	if !ok {
		return 0, false, errors.New("userId not found in token")
	}

	isAdmin, ok := claims["isAdmin"].(bool)
	if !ok {
		return 0, false, errors.New("isAdmin not found in token")
	}

	return uint(userIDFloat), isAdmin, nil
}

func GetUserFromHeader(c echo.Context) (uint, bool, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return 0, false, errors.New("Authorization header missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, false, errors.New("Invalid Authorization format")
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte("RaghavSureka"), nil
	})
	if err != nil || !token.Valid {
		return 0, false, errors.New("Invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false, errors.New("Invalid token claims")
	}

	userIDFloat, ok := claims["userId"].(float64)
	if !ok {
		return 0, false, errors.New("userId missing")
	}

	isAdmin, ok := claims["isAdmin"].(bool)
	if !ok {
		return 0, false, errors.New("isAdmin missing")
	}

	return uint(userIDFloat), isAdmin, nil
}
