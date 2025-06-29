// handler/medicine_handler.go
package handlers

import (
	"io"
	"net/http"
	"os"
	"pharmacy/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// MedicineHandler holds the methods for CRUD operations on Medicine
type MedicineHandler struct{}

// NewMedicineHandler creates a new instance of the handler
func NewMedicineHandler() *MedicineHandler {
	return &MedicineHandler{}
}

// CreateMedicine handles POST requests to create a new medicine
func (h *MedicineHandler) CreateMedicine(c echo.Context) error {
	var medicine models.Medicine
	if err := c.Bind(&medicine); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
	}

	// Call the model's CreateMedicine function to insert the medicine
	createdMedicine, err := models.CreateMedicine(medicine.Name, medicine.Description, medicine.CompanyID, medicine.UpdatedBy, medicine.Offer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not create medicine"})
	}

	return c.JSON(http.StatusCreated, createdMedicine)
}

// UpdateMedicine handles PUT requests to update an existing medicine
func (h *MedicineHandler) UpdateMedicine(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid medicine ID"})
	}

	var medicine models.Medicine
	if err := c.Bind(&medicine); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
	}

	// Call the model's UpdateMedicine function to update the medicine
	updatedMedicine, err := models.UpdateMedicine(uint(id), medicine.Name, medicine.Description, medicine.CompanyID, medicine.UpdatedBy, medicine.Offer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not update medicine"})
	}

	return c.JSON(http.StatusOK, updatedMedicine)
}

// DeleteMedicine handles DELETE requests to delete a medicine
func (h *MedicineHandler) DeleteMedicine(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid medicine ID"})
	}

	// Call the model's DeleteMedicine function to delete the medicine
	if err := models.DeleteMedicine(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not delete medicine"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Medicine deleted successfully"})
}

// GetMedicine handles GET requests to retrieve a specific medicine by ID
func (h *MedicineHandler) GetMedicine(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid medicine ID"})
	}

	medicine, err := models.GetMedicine(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not retrieve medicine"})
	}

	return c.JSON(http.StatusOK, medicine)
}

// GetAllMedicines handles GET requests to retrieve all medicines
func (h *MedicineHandler) GetAllMedicines(c echo.Context) error {
	medicines, err := models.GetAllMedicines()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not retrieve medicines"})
	}

	return c.JSON(http.StatusOK, medicines)
}

func (h *MedicineHandler) UploadMedicinesCSV(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "CSV file is required"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cannot open uploaded file"})
	}
	defer src.Close()

	tempPath := "/tmp/" + file.Filename
	dst, err := os.Create(tempPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cannot save file"})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cannot write file"})
	}

	// Pass to model logic
	if err := models.InsertMedicinesFromCSV(tempPath, "admin_user"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Medicines uploaded successfully"})
}

// UpdateOffer handles PUT requests to update offers for medicines
func (h *MedicineHandler) UpdateOffer(c echo.Context) error {
	var request struct {
		MedicineID uint   `json:"medicine_id"`
		CompanyID  uint   `json:"company_id"`
		Offer      string `json:"offer"`
		UpdatedBy  string `json:"updated_by"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
	}

	// Validate required fields
	if request.CompanyID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Company ID is required"})
	}

	if request.Offer == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Offer is required"})
	}

	if request.UpdatedBy == "" {
		request.UpdatedBy = "system" // Default value
	}

	// Call the model function to update the offer
	if err := models.UpdateOfferForMedicine(request.MedicineID, request.CompanyID, request.Offer, request.UpdatedBy); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not update offer"})
	}

	var message string
	if request.MedicineID == 0 {
		message = "Offer updated for all medicines in the company"
	} else {
		message = "Offer updated for the specific medicine"
	}

	return c.JSON(http.StatusOK, map[string]string{"message": message})
}
