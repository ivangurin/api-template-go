package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ivangurin/restful-api-go/database"
	"github.com/ivangurin/restful-api-go/models"
	"github.com/ivangurin/restful-api-go/routes"
)

func main() {

	if err := database.Connect(); err != nil {
		log.Fatal(err.Error())
	}

	database.AutoMigrate(&models.Document{}, &models.DocumentItem{})

	app := fiber.New()

	app.Use(cors.New())

	// Root
	app.Get("/", routes.GetRoot)

	// Documents
	app.Get("api/documents", routes.GetDocuments)
	app.Post("api/documents", routes.CreateDocument)
	app.Get("api/documents/:id", routes.GetDocument)
	app.Put("api/documents/:id", routes.UpdateDocument)
	app.Delete("api/documents/:id", routes.DeleteDocument)
	app.Post("api/documents/populate", routes.PopulateDocument)

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err.Error())
	}

}
