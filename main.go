package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/riskiapl/fiber-app/config"
	"github.com/riskiapl/fiber-app/cron"
	"github.com/riskiapl/fiber-app/database"
	"github.com/riskiapl/fiber-app/middleware"
	"github.com/riskiapl/fiber-app/routes"
	"github.com/riskiapl/fiber-app/utils"
)

func main() {
	// Load environment variables based on environment
	config.LoadEnv()

	// Log which environment we're running in
	env := config.GetEnvironment()
	log.Printf("Application running in %s mode\n", env)

	// Membuat database
	database.CreateDatabase()

	// Koneksi ke database
	database.ConnectDB()

	// Migrasi database
	database.Migrate()

	// Start cron jobs
	cron.StartCronJobs()

	// Initialize i18n
	if err := utils.InitGlobalI18n(); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}

	app := fiber.New()

	// Get CORS configuration based on environment
	corsConfig := config.GetCorsConfig()

	// Setup CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsConfig["origins"],
		AllowHeaders:     corsConfig["headers"],
		AllowMethods:     corsConfig["methods"],
		AllowCredentials: corsConfig["credentials"] == "true",
	}))

	// Add the i18n middleware to your fiber app
	app.Use(middleware.I18nMiddleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	// Setup routes
	routes.AuthRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
