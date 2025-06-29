package models

import (
	"time"
)

type Order struct {
	ID        uint        `json:"orderId" gorm:"primaryKey"`
	UserID    uint        `json:"userId"`
	Items     []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Status    string      `json:"status" gorm:"default:'pending'"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	UpdatedBy string      `json:"updated_by"`
	User      User        `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

// TableName specifies the table name for GORM to use
func (Order) TableName() string {
	return "order"
}

type OrderItem struct {
	ID         uint     `json:"-" gorm:"primaryKey"`
	OrderID    uint     `json:"-"`
	MedicineID uint     `json:"medicineId"`
	CompanyID  uint     `json:"companyId"`
	Quantity   int      `json:"quantity"`
	Medicine   Medicine `json:"medicine" gorm:"foreignKey:MedicineID;references:ID"`
	Company    Company  `json:"company" gorm:"foreignKey:CompanyID;references:ID"`
}

// TableName specifies the table name for GORM to use
func (OrderItem) TableName() string {
	return "order_item"
}

type OrderRequest struct {
	OrderID     uint               `json:"orderId"`
	UserID      uint               `json:"userId"`
	Items       []OrderItemRequest `json:"items"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"createdAt"`
	UserDetails *UserDetails       `json:"userDetails,omitempty"`
}

type UserDetails struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	FirmName string `json:"firmName"`
}

type OrderItemRequest struct {
	MedicineID   uint   `json:"medicineId"`
	MedicineName string `json:"medicineName,omitempty"`
	CompanyID    uint   `json:"companyId"`
	CompanyName  string `json:"companyName,omitempty"`
	Quantity     int    `json:"quantity"`
}

func ConvertOrderToOrderRequest(order *Order, includeUserDetails bool) *OrderRequest {
	var items []OrderItemRequest
	for _, item := range order.Items {
		items = append(items, OrderItemRequest{
			MedicineID:   item.MedicineID,
			MedicineName: item.Medicine.Name,
			CompanyID:    item.CompanyID,
			CompanyName:  item.Company.CompanyName,
			Quantity:     item.Quantity,
		})
	}

	orderRequest := &OrderRequest{
		OrderID:   order.ID,
		UserID:    order.UserID,
		Items:     items,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}

	// Include user details only if requested (for admin users)
	if includeUserDetails && order.User.ID != 0 {
		orderRequest.UserDetails = &UserDetails{
			Name:     order.User.Name,
			Phone:    order.User.Phone,
			FirmName: order.User.FirmName,
		}
	}

	return orderRequest
}

func CreateOrderWithItems(req OrderRequest, updatedBy string) (*OrderRequest, error) {
	order := Order{
		UserID:    req.UserID,
		Status:    "pending",
		UpdatedBy: updatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&order).Error; err != nil {
		return nil, err
	}

	var items []OrderItem
	for _, item := range req.Items {
		items = append(items, OrderItem{
			OrderID:    order.ID,
			MedicineID: item.MedicineID,
			CompanyID:  item.CompanyID,
			Quantity:   item.Quantity,
		})
	}
	if err := db.Create(&items).Error; err != nil {
		return nil, err
	}

	// Reload with associations
	db.Preload("Items.Medicine").Preload("Items.Company").First(&order)

	return ConvertOrderToOrderRequest(&order, false), nil
}

func GetOrderByID(id uint) (*OrderRequest, error) {
	var order Order
	if err := db.Preload("Items.Medicine").Preload("Items.Company").First(&order, id).Error; err != nil {
		return nil, err
	}
	return ConvertOrderToOrderRequest(&order, false), nil
}

func GetOrderByIDWithUserDetails(id uint, includeUserDetails bool) (*OrderRequest, error) {
	var order Order
	query := db.Preload("Items.Medicine").Preload("Items.Company")

	if includeUserDetails {
		query = query.Preload("User")
	}

	if err := query.First(&order, id).Error; err != nil {
		return nil, err
	}
	return ConvertOrderToOrderRequest(&order, includeUserDetails), nil
}

func GetAllOrders(userID uint, isAdmin bool) ([]OrderRequest, error) {
	var orders []Order
	query := db.Preload("Items.Medicine").Preload("Items.Company")

	// Preload user data only if admin is requesting
	if isAdmin {
		query = query.Preload("User")
	}

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	var result []OrderRequest
	for _, order := range orders {
		result = append(result, *ConvertOrderToOrderRequest(&order, isAdmin))
	}
	return result, nil
}

func UpdateOrder(id uint, req OrderRequest, updatedBy string) (*OrderRequest, error) {
	var order Order
	if err := db.First(&order, id).Error; err != nil {
		return nil, err
	}

	if err := db.Where("order_id = ?", id).Delete(&OrderItem{}).Error; err != nil {
		return nil, err
	}

	var items []OrderItem
	for _, item := range req.Items {
		items = append(items, OrderItem{
			OrderID:    id,
			MedicineID: item.MedicineID,
			CompanyID:  item.CompanyID,
			Quantity:   item.Quantity,
		})
	}
	if err := db.Create(&items).Error; err != nil {
		return nil, err
	}

	order.UserID = req.UserID
	order.UpdatedAt = time.Now()
	order.UpdatedBy = updatedBy
	if err := db.Save(&order).Error; err != nil {
		return nil, err
	}

	db.Preload("Items.Medicine").Preload("Items.Company").First(&order)

	order.Items = items
	return ConvertOrderToOrderRequest(&order, false), nil
}

func UpdateOrderStatus(id uint, status string, updatedBy string) (*OrderRequest, error) {
	var order Order
	if err := db.First(&order, id).Error; err != nil {
		return nil, err
	}

	order.Status = status
	order.UpdatedBy = updatedBy
	order.UpdatedAt = time.Now()

	if err := db.Save(&order).Error; err != nil {
		return nil, err
	}

	// Reload with associations
	if err := db.Preload("Items.Medicine").Preload("Items.Company").First(&order, id).Error; err != nil {
		return nil, err
	}

	return ConvertOrderToOrderRequest(&order, false), nil
}

func DeleteOrder(id uint) error {
	var order Order
	if err := db.First(&order, id).Error; err != nil {
		return err
	}
	if err := db.Delete(&OrderItem{}, "order_id = ?", id).Error; err != nil {
		return err
	}
	return db.Delete(&order).Error
}
