// jwt_util.go

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"go-server/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret string

func init() {
	GetSecretKey()
}

// Claims represents the JWT claims.
type Claims struct {
	UserId    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`
	jwt.StandardClaims
}

func GetSecretKey() {
	jwtSecret = os.Getenv("JWT_SECRET")
}

func GenerateJWTToken(user *models.User) (string, error) {
	claims := Claims{
		UserId:    int(user.ID),
		UserEmail: user.Email,
		UserRole:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

func ExtractClaims(r *http.Request) (jwt.MapClaims, bool) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, false
	}
	tokenString = tokenString[len("Bearer "):]

	token, err := ValidateJWTToken(tokenString)
	if err != nil || !token.Valid {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}

func VerifyJWTToken(tokenString string) (Claims, bool) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return claims, false
	}
	return claims, token.Valid
}

// GeneratePasswordResetToken generates a password reset token for a user.
func GeneratePasswordResetToken(user *models.User) (string, error) {
	// Create a unique token based on user email and a timestamp
	tokenData := user.Email + time.Now().String()

	// Generate a random 32-byte string for additional security
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Combine the token data and random bytes
	combinedData := append([]byte(tokenData), randomBytes...)

	// Encode the combined data as a base64 string
	token := base64.URLEncoding.EncodeToString(combinedData)

	return token, nil
}

// isTokenExpired checks if a reset token has expired.
func IsTokenExpired(tokenExpiry time.Time) bool {
	// You can define an expiration duration, e.g., 1 hour.
	expirationDuration := 1 * time.Hour

	// Calculate the expiration time by subtracting the duration from the current time.
	expirationTime := time.Now().Add(-expirationDuration)

	// Compare the token's expiry time with the calculated expiration time.
	// If the token's expiry time is before the calculated time, it's expired.
	return tokenExpiry.Before(expirationTime)
}
