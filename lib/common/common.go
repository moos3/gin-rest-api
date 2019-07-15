package common

import (
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// JSON - simple map keyed by string to a interface
type JSON = map[string]interface{}

var gdprCountries = []string{"AT", "BE", "HR", "BG", "CY", "CZ", "DK", "EE", "FI", "FR", "DE", "GR", "HU", "IE",
	"IT", "LV", "LT", "LU", "MT", "NL", "PL", "PT", "RO", "SK", "SI", "ES", "SE", "GB"}

// Hash - password hashing
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckHash - checking the password hash
func CheckHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJwtToken - Gerneate validate JWT token from common.JSON
func GenerateJwtToken(data JSON) (string, error) {
	region := os.Getenv("REGION")
	date := time.Now().Add(time.Hour * 24 * 7)
	username := data["username"]
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":        data,
		"exp":         date.Unix(),
		"host_region": region,
		"iss":         username,
	})

	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"
	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		return "", readErr
	}
	tokenString, err := token.SignedString(key)
	return tokenString, err
}
