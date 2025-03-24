package controllers

import (
	"time"

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
			"error": utils.T("auth.invalid_request_body"),
		})
	}

	// Validasi input
	if input.Userormail == "" || input.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.username_password_required"),
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
			"error": utils.T("auth.token_generation_error"),
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
			"error": utils.T("auth.invalid_request_body"),
		})
	}

	// Validasi input
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.registration_fields_required"),
		})
	}

	// Set plain password
	input.PlainPassword = input.Password

	// Proses register menggunakan service
	result, err := c.authService.Register(input)
	if err != nil {
		// Get the error code if it's an AppError
		errorCode := utils.GetErrorCode(err)

		// Check if it's a conflict error
		if errorCode == utils.ErrUsernameExists || errorCode == utils.ErrEmailExists {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func (c *AuthController) VerifyOTP(ctx *fiber.Ctx) error {
	var input types.VerifyOTPInput

	// Parse request body
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.invalid_request_body"),
		})
	}

	// Process OTP verification
	err := c.authService.VerifyRegistration(input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": utils.T("auth.email_verification_successful"),
	})
}

func (c *AuthController) CheckUsername(ctx *fiber.Ctx) error {
	// Get the username from query parameters
	username := ctx.Query("username")
	if username == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.username_parameter_required"),
		})
	}

	// Check if username exists using the auth service
	exists, err := c.authService.IsUsernameExists(username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": utils.T("auth.username_check_error"),
		})
	}

	// Return status 200 in both cases, just change the available flag
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"available": !exists,
		"username":  username,
	})
}

func (c *AuthController) DeletePendingMember(ctx *fiber.Ctx) error {
	// Get email from query parameter
	email := ctx.Query("email")
	if email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.email_parameter_required"),
		})
	}

	// Process delete pending member
	err := c.authService.DeletePendingMember(email)
	if err != nil {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": utils.T("auth.pending_member_deleted"),
	})
}

func (c *AuthController) ResendOTP(ctx *fiber.Ctx) error {
	var input types.ResendOTPInput

	// Parse request body
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.invalid_request_body"),
		})
	}

	// Validate the email from body
	if input.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.email_required"),
		})
	}

	// Process resend OTP
	result, err := c.authService.ResendOTP(input.Email)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (c *AuthController) ResetPassword(ctx *fiber.Ctx) error {
	var input types.ForgotPasswordInput

	// Parse request body
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.invalid_request_body"),
		})
	}

	// Validate the email from body
	if input.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.email_required"),
		})
	}

	// Process forgot password
	result, err := c.authService.ResetPassword(input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (c *AuthController) ChangePassword(ctx *fiber.Ctx) error {
	var input types.ChangePasswordInput

	// Parse request body
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.invalid_request_body"),
		})
	}

	// Validate input
	if input.Token == "" || input.NewPassword == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.T("auth.token_password_required"),
		})
	}

	// Parse and validate JWT token
	claims, err := utils.ParseToken(input.Token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": utils.T("auth.invalid_expired_token"),
		})
	}

	// Create ResetPasswordData from token claims
	resetPasswordData := types.ResetPasswordData{
		Email: claims["email"].(string),
		Token: claims["token"].(string),
	}

	// Convert expires_at from Unix timestamp to time.Time
	if expired, ok := claims["expired"].(float64); ok {
		resetPasswordData.Expired = time.Unix(int64(expired), 0)
	}

	// Pass token claims and new password to change password function
	err = c.authService.ChangePassword(resetPasswordData, input.NewPassword)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": utils.T("auth.password_changed_successfully"),
	})
}
