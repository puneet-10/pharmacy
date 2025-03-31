// handler/medicine_handler.go
package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"pharmacy/models"
	"strconv"
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
	createdMedicine, err := models.CreateMedicine(medicine.Name, medicine.Description, medicine.CompanyID, medicine.UpdatedBy)
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
	updatedMedicine, err := models.UpdateMedicine(uint(id), medicine.Name, medicine.Description, medicine.CompanyID, medicine.UpdatedBy)
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
