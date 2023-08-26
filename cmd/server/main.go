package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jackramey/totp/internal/db"
	"github.com/jackramey/totp/totp"
)

var generator = totp.NewGenerator(totp.WithX(30), totp.WithT0(0), totp.WithD(6))

func main() {
	dbClient := openDb()
	defer dbClient.Close()

	// Initialize the database
	if err := db.InitDatabase(dbClient); err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	e.POST("/register", RegisterHandler)
	e.POST("/verify", VerifyHandler)
	e.POST("/challenge", ChallengeHandler)

	// Start the server
	e.Start(":3000")
}

type RegistrationPayload struct {
	Username string `json:"username"`
}

type VerificationPayload struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

func RegisterHandler(c echo.Context) error {
	dbClient := openDb()
	defer dbClient.Close()

	payload := new(RegistrationPayload)
	if err := c.Bind(payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON payload")
	}

	// Generate a new secret for the user
	secret := generateSecret()

	// Store the user's information in the database
	if err := db.InsertUser(dbClient, payload.Username, secret); err != nil {
		return c.String(http.StatusInternalServerError, "Error inserting user data")
	}

	return c.String(http.StatusCreated, secret)
}

func VerifyHandler(c echo.Context) error {
	dbClient := openDb()
	defer dbClient.Close()

	payload := new(VerificationPayload)
	if err := c.Bind(payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON payload")
	}

	code, err := strconv.Atoi(payload.Code)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid code")
	}

	// Retrieve user's secret from the database
	secret, err := db.GetSecretByUsername(dbClient, payload.Username)
	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}

	// Validate the code using totp.Generate from your package
	generatedCode, _, err := generator.Generate(secret)
	if err != nil || uint32(code) != generatedCode {
		return c.String(http.StatusUnauthorized, "Invalid code")
	}

	// Update the user's verification status in the database
	if err := db.UpdateVerificationStatus(dbClient, payload.Username); err != nil {
		return c.String(http.StatusInternalServerError, "Error updating verification status")
	}

	return c.String(http.StatusOK, "Verification successful")
}

func ChallengeHandler(c echo.Context) error {
	dbClient := openDb()
	defer dbClient.Close()

	payload := new(VerificationPayload)
	if err := c.Bind(payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON payload")
	}

	code, err := strconv.Atoi(payload.Code)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid code")
	}

	// Retrieve user's secret from the database
	secret, err := db.GetSecretByUsername(dbClient, payload.Username)
	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}

	// Validate the code using totp.Generate from your package
	generatedCode, _, err := generator.Generate(secret)
	if err != nil || uint32(code) != generatedCode {
		return c.String(http.StatusUnauthorized, "Invalid code")
	}

	return c.String(http.StatusOK, "Challenge successful")
}

func openDb() *sql.DB {
	dbClient, err := sql.Open("sqlite3", "server.db")
	if err != nil {
		panic(err)
	}

	return dbClient
}

func generateSecret() string {
	// Generate a random 20-byte secret
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		fmt.Println("Error generating secret:", err)
		return ""
	}

	// Encode the secret in base32
	encodedSecret := base32.StdEncoding.EncodeToString(secret)

	return encodedSecret
}
