package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riskiapl/fiber-app/services"
	"github.com/riskiapl/fiber-app/types"
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

	return ctx.Status(fiber.StatusOK).JSON(result)
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
