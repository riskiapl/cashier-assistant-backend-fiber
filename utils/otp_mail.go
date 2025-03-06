package utils

import (
	"fmt"
)

// OtpMail generates an HTML email body for OTP verification
func OTPMail(otpCode string, expiresAt string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OTP Verification</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f9f9f9;
        color: #333;
        margin: 0;
        padding: 0;
      }
      .container {
        max-width: 600px;
        margin: 50px auto;
        background: #ffffff;
        padding: 20px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        text-align: center;
      }
      .header {
        font-size: 24px;
        font-weight: bold;
        color: #4CAF50;
        margin-bottom: 10px;
      }
      .otp-code {
        font-size: 32px;
        font-weight: bold;
        color: #333;
        background-color: #f3f3f3;
        display: inline-block;
        padding: 10px 20px;
        border-radius: 5px;
        margin-top: 20px;
      }
      .footer {
        margin-top: 30px;
        font-size: 14px;
        color: #888;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">OTP Verification</div>
      <p>Hi,</p>
      <p>Thank you for signing up. Please use the following One-Time Password (OTP) to verify your email address:</p>
      <div class="otp-code">%s</div>
      <p>This OTP is valid for %s.</p>
      <p>If you did not request this code, please ignore this message.</p>
      <div class="footer">
        &copy; 2025 Cashier Assistant. All rights reserved.
      </div>
    </div>
  </body>
  </html>`, otpCode, expiresAt)
}
