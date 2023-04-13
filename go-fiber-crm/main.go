package main

import (
	"github.com/elizraa/go-fiber-crm/database"
	"github.com/elizraa/go-fiber-crm/lead"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func initDatabase() {
	database.Connect()
	db = database.GetDB()
	db.AutoMigrate(&lead.Lead{})
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Get("/api/v1/leads", lead.GetLeads)
	app.Post("/api/v1/lead", lead.NewLead)
	app.Delete("/api/v1/:id", lead.DeleteLead)
}

func main() {
	app := fiber.New()
	initDatabase()
	setupRoutes(app)
	app.Listen(8080)

	defer db.Close()

}
