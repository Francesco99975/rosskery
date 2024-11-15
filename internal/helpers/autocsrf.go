package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetCsrfTokenFromServer(sessionID string) (string, error) {
	// Define the URL for the CSRF token endpoint
	host := os.Getenv("HOST")
	url := fmt.Sprintf("%s/csrf", host) // Adjust the URL as needed

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the custom authentication header
	req.Header.Set("X-Auth-Token", os.Getenv("AUTH_CSRF_KEY")) // Set the token from environment variable

	//Spoof sessionid
	req.Header.Set("Cookie", fmt.Sprintf("session_id=%s", sessionID))

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response to get the CSRF token
	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Extract and return the CSRF token
	token, ok := result["csrfToken"]
	if !ok {
		return "", fmt.Errorf("csrfToken not found in response")
	}

	return token, nil
}
