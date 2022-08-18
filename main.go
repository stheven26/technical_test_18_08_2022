package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stheven26/db"
	"github.com/stheven26/routes"
)

func main() {
	db.SetupDB()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Router(app)

	app.Listen(":8000")
}
