package main

import (
	"log"
	"pharmacy/handlers"
	"pharmacy/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Database connection string (adjust this to your actual DB setup)
	connStr := "postgresql://pharamcy_owner:CSyuQK3I9WXl@ep-shrill-haze-a1zsvdfj.ap-southeast-1.aws.neon.tech/pharamcy?sslmode=require"
	// Connect to PostgreSQL using GORM
	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Auto-migrate the schema (creates the users table, and ensures it has the correct columns)
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Error migrating database:", err)
	}

	// Set the database connection for models
	models.SetDB(db)
	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		SigningKey:  []byte("RaghavSureka"),
		ContextKey:  "token",
	})

	// Initialize Echo instance
	e := echo.New()

	// Middleware: logging and recovery
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	companyHandler := handlers.NewCompanyHandler()
	medicineHandler := handlers.NewMedicineHandler()
	orderHandler := handlers.NewOrderHandler()

	// Define routes
	e.POST("/signup", handlers.SignUpHandler)
	e.POST("/authenticate", handlers.AuthenticateHandler)
	e.PUT("/user/:id", handlers.UpdateUserHandler)

	e.POST("/companies", companyHandler.CreateCompany)
	e.GET("/companies", companyHandler.GetAllCompanies)
	e.GET("/companies/:id", companyHandler.GetCompany)
	e.PUT("/companies/:id", companyHandler.UpdateCompany)
	e.DELETE("/companies/:id", companyHandler.DeleteCompany)

	e.POST("/medicines", medicineHandler.CreateMedicine)
	e.PUT("/medicines/:id", medicineHandler.UpdateMedicine)
	e.DELETE("/medicines/:id", medicineHandler.DeleteMedicine)
	e.GET("/medicines/:id", medicineHandler.GetMedicine)
	e.GET("/medicines", medicineHandler.GetAllMedicines)
	e.POST("/medicines/upload", medicineHandler.UploadMedicinesCSV)

	e.POST("/orders", orderHandler.CreateOrder)
	e.GET("/orders/:id", orderHandler.GetOrder)
	e.PUT("/orders/:id", orderHandler.UpdateOrder)
	e.DELETE("/orders/:id", orderHandler.DeleteOrder)
	e.GET("/orders", orderHandler.GetAllOrders, jwtMiddleware)

	e.PUT("/medicines/offer", medicineHandler.UpdateOffer)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
