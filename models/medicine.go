// models/medicine.go
package models

import (
	"time"
)

// Medicine struct represents the medicine model in the database
type Medicine struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CompanyID   uint      `json:"company_id"` // Foreign key for Company
	Company     Company   `json:"company"`    // Foreign key relationship with Company
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

type MedicineDTO struct {
	MedicineID uint   `json:"medicineId"`
	Name       string `json:"name"`
}

type CompanyMedicinesResponse struct {
	CompanyID   uint          `json:"companyId"`
	CompanyName string        `json:"companyName"`
	Medicines   []MedicineDTO `json:"medicines"`
}

// TableName overrides the default table name
func (Medicine) TableName() string {
	return "medicine" // This ensures GORM uses the singular "medicine" table name
}

// CreateMedicine creates a new medicine in the database
func CreateMedicine(name, description string, companyID uint, updatedBy string) (*Medicine, error) {
	medicine := &Medicine{
		Name:        name,
		Description: description,
		CompanyID:   companyID,
		UpdatedBy:   updatedBy,
	}

	// Insert the medicine into the database
	if err := db.Create(medicine).Error; err != nil {
		return nil, err
	}

	// Return the created medicine
	return medicine, nil
}

// UpdateMedicine updates an existing medicine in the database
func UpdateMedicine(id uint, name, description string, companyID uint, updatedBy string) (*Medicine, error) {
	var medicine Medicine

	// Find the medicine by ID
	if err := db.First(&medicine, id).Error; err != nil {
		return nil, err
	}

	// Update the fields
	medicine.Name = name
	medicine.Description = description
	medicine.CompanyID = companyID
	medicine.UpdatedBy = updatedBy
	medicine.UpdatedAt = time.Now()

	// Save the updated medicine to the database
	if err := db.Save(&medicine).Error; err != nil {
		return nil, err
	}

	return &medicine, nil
}

// DeleteMedicine deletes a medicine from the database
func DeleteMedicine(id uint) error {
	var medicine Medicine

	// Find the medicine by ID
	if err := db.First(&medicine, id).Error; err != nil {
		return err
	}

	// Delete the medicine
	if err := db.Delete(&medicine).Error; err != nil {
		return err
	}

	return nil
}

// GetMedicine retrieves a specific medicine by its ID
func GetMedicine(id uint) (*Medicine, error) {
	var medicine Medicine

	// Fetch the medicine by ID and preload its associated company data
	if err := db.Preload("Company").First(&medicine, id).Error; err != nil {
		return nil, err
	}

	return &medicine, nil
}

// GetAllMedicines retrieves all medicines and groups them by their associated company
//func GetAllMedicines() ([]Medicine, error) {
//	var medicines []Medicine
//
//	// Fetch all medicines with company details preloaded
//	if err := db.Preload("Company").Find(&medicines).Error; err != nil {
//		return nil, err
//	}
//
//	return medicines, nil
//}

func GetAllMedicines() ([]CompanyMedicinesResponse, error) {
	var medicines []Medicine

	// Fetch medicines with associated company data
	if err := db.Preload("Company").Find(&medicines).Error; err != nil {
		return nil, err
	}

	companyMap := make(map[uint]*CompanyMedicinesResponse)

	for _, med := range medicines {
		comp := med.Company
		if _, exists := companyMap[comp.ID]; !exists {
			companyMap[comp.ID] = &CompanyMedicinesResponse{
				CompanyID:   comp.ID,
				CompanyName: comp.CompanyName,
				Medicines:   []MedicineDTO{},
			}
		}

		medicineDTO := MedicineDTO{
			MedicineID: med.ID,
			Name:       med.Name,
		}

		companyMap[comp.ID].Medicines = append(companyMap[comp.ID].Medicines, medicineDTO)
	}

	// Convert map to slice
	var response []CompanyMedicinesResponse
	for _, v := range companyMap {
		response = append(response, *v)
	}

	return response, nil
}
