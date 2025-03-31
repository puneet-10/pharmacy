package models

import (
	"time"
)

// Company struct represents the company model in the database
type Company struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CompanyName string    `json:"company_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

// TableName specifies the table name for GORM to use
func (Company) TableName() string {
	return "company" // Use the singular name of your table
}

// CreateCompany creates a new company in the database
func CreateCompany(companyName, description, updatedBy string) (*Company, error) {
	company := &Company{
		CompanyName: companyName,
		Description: description,
		UpdatedBy:   updatedBy,
	}

	// Insert the company into the database
	if err := db.Create(company).Error; err != nil {
		return nil, err
	}

	// Return the created company
	return company, nil
}

// GetCompany retrieves a company by ID from the database
func GetCompany(id uint) (*Company, error) {
	var company Company
	if err := db.First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

// UpdateCompany updates an existing company in the database
func UpdateCompany(id uint, companyName, description, updatedBy string) (*Company, error) {
	var company Company

	// Find the company by ID
	if err := db.First(&company, id).Error; err != nil {
		return nil, err
	}

	// Update the company fields
	company.CompanyName = companyName
	company.Description = description
	company.UpdatedBy = updatedBy
	company.UpdatedAt = time.Now()

	// Save the updated company to the database
	if err := db.Save(&company).Error; err != nil {
		return nil, err
	}

	return &company, nil
}

// DeleteCompany deletes a company from the database by ID
func DeleteCompany(id uint) error {
	if err := db.Delete(&Company{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAllCompanies retrieves all companies from the database
func GetAllCompanies() ([]Company, error) {
	var companies []Company
	if err := db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}
