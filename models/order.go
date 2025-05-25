// models/order.go
package models

import (
	"time"
)

type Order struct {
	ID        uint        `json:"orderId" gorm:"primaryKey"`
	UserID    uint        `json:"userId"`
	Items     []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	UpdatedBy string      `json:"updated_by"`
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
	OrderID uint               `json:"orderId"`
	UserID  uint               `json:"userId"`
	Items   []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	MedicineID   uint   `json:"medicineId"`
	MedicineName string `json:"medicineName,omitempty"`
	CompanyID    uint   `json:"companyId"`
	CompanyName  string `json:"companyName,omitempty"`
	Quantity     int    `json:"quantity"`
}

func ConvertOrderToOrderRequest(order *Order) *OrderRequest {
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
	return &OrderRequest{
		OrderID: order.ID,
		UserID:  order.UserID,
		Items:   items,
	}
}

func CreateOrderWithItems(req OrderRequest, updatedBy string) (*OrderRequest, error) {
	order := Order{
		UserID:    req.UserID,
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

	return ConvertOrderToOrderRequest(&order), nil
}

func GetOrderByID(id uint) (*OrderRequest, error) {
	var order Order
	if err := db.Preload("Items.Medicine").Preload("Items.Company").First(&order, id).Error; err != nil {
		return nil, err
	}
	return ConvertOrderToOrderRequest(&order), nil
}

func GetAllOrders(userID uint, isAdmin bool) ([]OrderRequest, error) {
	var orders []Order
	query := db.Preload("Items.Medicine").Preload("Items.Company")

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	var result []OrderRequest
	for _, order := range orders {
		result = append(result, *ConvertOrderToOrderRequest(&order))
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
	return ConvertOrderToOrderRequest(&order), nil
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
