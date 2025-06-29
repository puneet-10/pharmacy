// models/medicine.go
package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
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
	Offer       string    `json:"offer"`
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
func CreateMedicine(name, description string, companyID uint, updatedBy string, offer string) (*Medicine, error) {
	medicine := &Medicine{
		Name:        name,
		Description: description,
		CompanyID:   companyID,
		UpdatedBy:   updatedBy,
		Offer:       offer,
	}

	// Insert the medicine into the database
	if err := db.Create(medicine).Error; err != nil {
		return nil, err
	}

	// Return the created medicine
	return medicine, nil
}

// UpdateMedicine updates an existing medicine in the database
func UpdateMedicine(id uint, name, description string, companyID uint, updatedBy string, offer string) (*Medicine, error) {
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
	medicine.Offer = offer
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

type ParsedMedicine struct {
	Name        string
	Description string
	CompanyName string
}

func InsertMedicinesFromCSV(filePath string, updatedBy string) error {
	// Step 1: Read and parse the CSV
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read() // skip header

	// Step 2: Load all companies
	var companies []Company
	if err := db.Find(&companies).Error; err != nil {
		return err
	}
	companyMap := make(map[string]Company)
	for _, c := range companies {
		companyMap[strings.ToLower(strings.TrimSpace(c.CompanyName))] = c
	}

	// Step 3: Process and insert each record
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 3 {
			continue // skip bad lines
		}

		name := strings.TrimSpace(record[0])
		desc := strings.TrimSpace(record[1])
		companyName := strings.TrimSpace(record[2])
		lookup := strings.ToLower(companyName)

		company, exists := companyMap[lookup]
		if !exists {
			// Create new company
			company = Company{
				CompanyName: companyName,
				Description: "Auto-generated via CSV",
				UpdatedBy:   updatedBy,
			}
			if err := db.Create(&company).Error; err != nil {
				return fmt.Errorf("error inserting company: %w", err)
			}
			companyMap[lookup] = company
		}

		medicine := Medicine{
			Name:        name,
			Description: desc,
			CompanyID:   company.ID,
			UpdatedBy:   updatedBy,
		}
		if err := db.Create(&medicine).Error; err != nil {
			return fmt.Errorf("error inserting medicine: %w", err)
		}
	}

	return nil
}

// UpdateOfferForMedicine updates offer for a specific medicine or all medicines in a company
func UpdateOfferForMedicine(medicineID uint, companyID uint, offer string, updatedBy string) error {
	if medicineID == 0 {
		// Update offer for all medicines under the company
		if err := db.Model(&Medicine{}).Where("company_id = ?", companyID).Updates(map[string]interface{}{
			"offer":      offer,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}
	} else {
		// Update offer for specific medicine
		if err := db.Model(&Medicine{}).Where("id = ?", medicineID).Updates(map[string]interface{}{
			"offer":      offer,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
