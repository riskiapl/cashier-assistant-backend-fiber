package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/riskiapl/fiber-app/database"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	log.Fatal(app.Listen(":8000"))
}
