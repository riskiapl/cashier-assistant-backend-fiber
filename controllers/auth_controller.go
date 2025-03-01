package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riskiapl/fiber-app/services"
	"github.com/riskiapl/fiber-app/types"
	"github.com/riskiapl/fiber-app/utils"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var input types.LoginInput

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi input
	if input.Userormail == "" || input.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username or password are required",
		})
	}

	// Proses login menggunakan service
	result, err := c.authService.Login(input)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Create a map with a single "data" key containing the entire result
	tokenPayload := map[string]any{
		"data": result,
	}

	// Generate JWT token
	token, err := utils.GenerateToken(tokenPayload)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating token",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
		"user":  result,
	})
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var input types.RegisterInput

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi input
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email and password are required",
		})
	}

	// Set plain password
	input.PlainPassword = input.Password

	// Proses register menggunakan service
	result, err := c.authService.Register(input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func (c *AuthController) CheckUsername(ctx *fiber.Ctx) error {
	// Get the username from query parameters
	username := ctx.Query("username")
	if username == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username parameter is required",
		})
	}

	// Check if username exists using the auth service
	exists, err := c.authService.IsUsernameExists(username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking username availability",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"available": !exists,
		"username":  username,
	})
}
