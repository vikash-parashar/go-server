// handlers.go

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vikash-parashar/asset-locator/db"
	"github.com/vikash-parashar/asset-locator/logger"
	"github.com/vikash-parashar/asset-locator/models"
	"github.com/vikash-parashar/asset-locator/utils"

	"github.com/gin-gonic/gin"
)

// SignUp handles the registration of a new user.
func SignUp(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InfoLogger.Println("Handling POST request for user registration")

		var signupRequest struct {
			FirstName string `json:"first_name" binding:"required"`
			LastName  string `json:"last_name" binding:"required"`
			Phone     string `json:"phone" binding:"required"`
			Email     string `json:"email" binding:"required"`
			Password  string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&signupRequest); err != nil {
			logger.ErrorLogger.Println("Invalid form data for user registration:", err)
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid form data"})
			return
		}

		// Check if the user already exists (by email or any other unique identifier)
		_, err := db.GetUserByEmailID(signupRequest.Email)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "User with this email already exists"})
			return
		}
		// Create a new user
		newUser := &models.User{
			FirstName: signupRequest.FirstName,
			LastName:  signupRequest.LastName,
			Phone:     signupRequest.Phone,
			Email:     signupRequest.Email,
			Password:  signupRequest.Password,
		}

		if newUser.Email == "gowithvikash@gmail.com" {
			newUser.Role = "admin"
		} else {
			newUser.Role = "general"
		}
		// Hash the password
		hashedPassword, err := utils.HashPassword(newUser.Password)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to hash password"})
			return
		}
		newUser.Password = hashedPassword

		if err := db.RegisterUser(newUser); err != nil {
			logger.ErrorLogger.Println("Failed to create user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user"})
			return
		}

		logger.InfoLogger.Println("User registered successfully")
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "User registered successfully"})
	}
}

// Login handles the user login and returns a JWT token upon successful login.
func Login(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InfoLogger.Println("Handling POST request for user login")

		var loginRequest struct {
			Email    string `form:"email" binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		if err := c.ShouldBind(&loginRequest); err != nil {
			logger.ErrorLogger.Println("Invalid form data for user login:", err)
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid form data"})
			return
		}

		fmt.Println(loginRequest)

		// Check if the user exists in the database
		user, err := db.GetUserByEmailID(loginRequest.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Incorrect email or password"})

			return
		}

		// Verify the password
		if !utils.VerifyPassword(loginRequest.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Incorrect password"})
			return
		}

		// Generate a JWT token
		token, err := utils.GenerateJWTToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate JWT token"})
			return
		}

		cookie := http.Cookie{
			Name:    "jwt-token",
			Value:   token,
			Expires: time.Now().Add(60 * time.Minute),
		}
		http.SetCookie(c.Writer, &cookie)

		logger.InfoLogger.Println("User logged in successfully")
		c.JSON(http.StatusOK, gin.H{"success": true, "token": token, "message": "Login successful"})
	}
}

// Logout handles the user logout by clearing the JWT token cookie.
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InfoLogger.Println("Handling GET request for user logout")

		// Clear the JWT token cookie by setting its expiration to a past time
		cookie := http.Cookie{
			Name:     "jwt-token",
			Value:    "",              // Clear the cookie value
			Expires:  time.Unix(0, 0), // Set the expiration time to Unix epoch (1970-01-01 00:00:00 UTC)
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Secure:   true, // Set this to true if you're using HTTPS
		}
		http.SetCookie(c.Writer, &cookie)
		c.Redirect(http.StatusPermanentRedirect, "/")
		logger.InfoLogger.Println("User logged out successfully")
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout successful"})
	}
}

// ForgotPassword handles the process of resetting a user's forgotten password.
func ForgotPassword(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InfoLogger.Println("Handling POST request for password reset")

		// Retrieve email address from the user input
		var resetRequest struct {
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&resetRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input data"})
			return
		}

		// Check if the user exists in the database
		user, err := db.GetUserByEmailID(resetRequest.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
			return
		}

		// Generate a unique reset token and set an expiration time for it (e.g., 1 hour)
		resetToken, err := utils.GeneratePasswordResetToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate reset token"})
			return
		}

		expiryTime := time.Now().Add(1 * time.Hour)
		// Save the reset token in the database associated with the user's account
		if err := db.SetResetToken(int(user.ID), resetToken, expiryTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save reset token"})
			return
		}

		// Send an email to the user with the reset URL
		err = utils.SendResetPasswordEmail(user.Email, resetToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send reset email"})
			return
		}

		logger.InfoLogger.Println("Password reset instructions sent successfully")
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Reset instructions sent to your email"})
	}
}

func ResetPassword(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InfoLogger.Println("Handling POST request for resetting password")

		resetToken := c.Query("token")

		if resetToken == "" {
			logger.ErrorLogger.Println("Reset token is missing")
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reset token is missing"})
			return
		}

		log.Println(resetToken)
		// Parse the new password from the request body
		var resetRequest struct {
			NewPassword string `json:"new_password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&resetRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input data"})
			return
		}

		// Verify the reset token
		user, err := db.VerifyResetToken(resetToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired reset token"})
			return
		}

		// Hash the new password
		hashedPassword, err := utils.HashPassword(resetRequest.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash the new password"})
			return
		}

		// Update the user's password in the database
		if err := db.UpdateUserPassword(int(user.ID), hashedPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update the password"})
			return
		}

		// Clear the reset token from the database
		if err := db.ClearResetToken(int(user.ID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to clear the reset token"})
			return
		}

		logger.InfoLogger.Println("Password reset successful")
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Password reset successful"})
	}
}

func GetCurrentUser(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InfoLogger.Println("Handling GET request for current user details")

		// Retrieve the JWT token from the cookie
		cookie, err := c.Request.Cookie("jwt-token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		token := cookie.Value

		claims, valid := utils.VerifyJWTToken(token)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract the user email from the claims
		userEmail := claims.UserEmail

		// Retrieve the user based on the user email from the database
		user, err := db.GetUserByEmailID(userEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user"})
			c.Abort()
			return
		}

		// Send the user information in the response
		logger.InfoLogger.Println("Current user details retrieved successfully")
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func RenderResetPasswordPage(c *gin.Context) {
	logger.InfoLogger.Println("Rendering reset password page")
	c.HTML(http.StatusOK, "reset_password.html", gin.H{})
}
