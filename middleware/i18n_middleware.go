package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riskiapl/fiber-app/utils"
)

// I18nMiddleware sets the language for the current request based on Accept-Language header
func I18nMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Accept-Language header
		lang := c.Get("Accept-Language")

		// If language is empty, set default to 'en'
		if lang == "" {
			lang = "en"
		}

		// Set the language for the current request
		if utils.GlobalI18n != nil {
			utils.GlobalI18n.SetLanguage(lang)
		}

		// Store the language in the context locals for later use if needed
		c.Locals("language", lang)

		return c.Next()
	}
}
