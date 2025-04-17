package agg_tcp

import (
	"CloudlogAutoLogger/internal/agg_logger"
	"bytes"
	"encoding/json"
	"net/http"
)

func SendADI2Cloudlog(url string, payload map[string]string) bool {

	// Convert payload to JSON
	_jsonPayload, err := MapToJSONStringNoEscape(payload)
	if err != nil {
		agg_logger.Get().Log("Error marshaling JSON:", err.Error())
		return false
	}
	_jsonBytes := []byte(_jsonPayload)

	// Create the HTTP request
	_req, err := http.NewRequest("POST", url, bytes.NewBuffer(_jsonBytes))
	if err != nil {
		agg_logger.Get().Log("Error creating request:", err.Error())
		return false
	}

	// Set the content type to application/json
	_req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(_req)
	if err != nil {
		agg_logger.Get().Log("Error making request:", err.Error())
		return false
	}
	defer resp.Body.Close()

	return true
}

func MapToJSONStringNoEscape(data map[string]string) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // Disable escaping of HTML characters

	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	// Remove the trailing newline added by encoder.Encode
	return buf.String()[:buf.Len()-1], nil
}
