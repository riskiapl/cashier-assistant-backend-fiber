package types

import "time"

// LoginInput represents the input data for login
type LoginInput struct {
	Userormail string `json:"userormail" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

// LoginResponse represents the response data for login
type LoginResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
}

// RegisterInput represents the input data for registration
type RegisterInput struct {
	Username      string `json:"username" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=6"`
	PlainPassword string `json:"plain_password"`
}

// RegisterResponse represents the response data for registration
type RegisterResponse struct {
	Success   string    `json:"success"`
	ExpiredAt time.Time `json:"expired_at"`
}

// VerifyOTPInput represents the input data for OTP verification
type VerifyOTPInput struct {
	Email   string `json:"email" validate:"required,email"`
	OtpCode string `json:"otp_code" validate:"required"`
}

// VerifyOTPResponse represents the response data for OTP verification
type VerifyOTPResponse struct {
	Message string `json:"message"`
}

type ResendOTPInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}
