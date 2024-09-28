package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TokenAuthentication is a middleware for token validation.
func TokenAuthentication(c *fiber.Ctx) error {
	// Extract the token from the "Authorization" header
	token := c.Get("Authorization")

	// Check if the token is empty or does not start with "Bearer "
	if token == "" || token[:7] != "Bearer " {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "خطا در اعتبار سنجی. کد خطا 55",
		})
	}

	// Extract the token string (remove "Bearer " prefix)
	token = token[7:]
	validationURL := "https://estate.sedrehgroup.ir/api/meets/verify/token/"

	// Create a map for the request body
	requestBody := map[string]string{
		"token": token,
	}

	// Marshal the map into JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "خطای داخلی رخ داده است. کد خطا 52",
		})
	}

	// Create a request with the JSON body
	req, err := http.NewRequest("POST", validationURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "خطای داخلی رخ داده است. کد خطا 52",
		})
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // Set timeout to 10 seconds
	}

	// Make the request to the API using the custom client
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "خطای داخلی رخ داده است. کد خطا 53",
		})
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "توکن معتبر نمی باشد. کد خطا 54",
		})
	}

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "خطای داخلی رخ داده است. کد خطا 56",
		})
	}

	// Initialize a variable to hold the parsed response as a map
	var parsedResponse map[string]interface{}

	// Unmarshal the response body into the parsedResponse map
	if err := json.Unmarshal(responseBody, &parsedResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "خطای داخلی رخ داده است. کد خطا 57",
		})
	}
	c.Locals("user_id", parsedResponse["user_id"])
	// c.Locals("user_id", token)

	// Token is valid, continue with the next middleware
	return c.Next()
}
