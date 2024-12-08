package common

import (
	"fmt"
)

func getEmailContentFormat(otp string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>OTP Verification</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f4f4f4;
					color: #333333;
				}
				.container {
					max-width: 600px;
					margin: 20px auto;
					background-color: #ffffff;
					border: 1px solid #dddddd;
					border-radius: 8px;
					box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
					overflow: hidden;
				}
				.header {
					background-color: #4CAF50;
					padding: 20px;
					text-align: center;
					color: #ffffff;
					font-size: 24px;
				}
				.content {
					padding: 20px;
					text-align: center; /* Center all content within this section */
				}
				.content p {
					margin: 10px 0;
					line-height: 1.6;
				}
				.otp {
					display: inline-block;
					font-size: 24px;
					font-weight: bold;
					color: #4CAF50;
					background-color: #f9f9f9;
					padding: 10px 20px;
					margin: 20px 0;
					border-radius: 4px;
					border: 1px dashed #4CAF50;
					text-align: center; /* Ensure the OTP itself is centered */
				}
				.footer {
					background-color: #f9f9f9;
					padding: 15px;
					text-align: center;
					font-size: 12px;
					color: #777777;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<!-- Header Section -->
				<div class="header">
					MePass Verification
				</div>

				<!-- Content Section -->
				<div class="content">
					<p>Hello,</p>
					<p>We received a request to verify your account or reset your password. Use the OTP below to proceed:</p>
					<div class="otp">%s</div>
					<p>If you did not request this, please ignore this email or contact our support team if you have concerns.</p>
					<p>Thank you,<br>The MePass Team</p>
				</div>

				<!-- Footer Section -->
				<div class="footer">
					&copy; 2024 MePass. All rights reserved.
				</div>
			</div>
		</body>
		</html>
	`, otp)
}
