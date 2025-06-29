// handler/company_handler.go
package handlers

import (
	"net/http"
	"pharmacy/models"
	"strconv"
	"github.com/labstack/echo/v4"
)

// CompanyHandler holds the database connection and provides methods to handle HTTP requests
type CompanyHandler struct{}

// NewCompanyHandler creates a new instance of the handler
func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{}
}

// CreateCompany handles POST requests to create a new company
func (h *CompanyHandler) CreateCompany(c echo.Context) error {
	var company models.Company
	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
	}

	// Call the model's CreateCompany function to insert the company
	createdCompany, err := models.CreateCompany(company.CompanyName, company.Description, company.UpdatedBy, company.LogoUrl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not create company"})
	}

	return c.JSON(http.StatusCreated, createdCompany)
}

// GetCompany handles GET requests to retrieve a company by ID
func (h *CompanyHandler) GetCompany(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid company ID"})
	}

	company, err := models.GetCompany(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Company not found"})
	}

	return c.JSON(http.StatusOK, company)
}

// UpdateCompany handles PUT requests to update an existing company
func (h *CompanyHandler) UpdateCompany(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid company ID"})
	}

	var company models.Company
	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
	}

	// Call the model's UpdateCompany function to update the company
	updatedCompany, err := models.UpdateCompany(uint(id), company.CompanyName, company.Description, company.UpdatedBy, company.LogoUrl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not update company"})
	}

	return c.JSON(http.StatusOK, updatedCompany)
}

// DeleteCompany handles DELETE requests to remove a company
func (h *CompanyHandler) DeleteCompany(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid company ID"})
	}

	// Call the model's DeleteCompany function to delete the company
	if err := models.DeleteCompany(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not delete company"})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// GetAllCompanies handles GET requests to retrieve all companies
func (h *CompanyHandler) GetAllCompanies(c echo.Context) error {
	companies, err := models.GetAllCompanies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not retrieve companies"})
	}

	return c.JSON(http.StatusOK, companies)
}
