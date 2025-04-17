package agg_tcp

import (
	"CloudlogAutoLogger/internal/agg_logger"
	"bytes"
	"encoding/json"
	"net/http"
)

type Cloudlog_payload struct {
	Cloudlog_api_key   string
	Station_profile_id string
	Payload            string
}

func (cp *Cloudlog_payload) SendADI2Cloudlog() bool {
	// Convert payload to string map
	_payload := make(map[string]string)

	_payload["key"] = cp.Cloudlog_api_key
	_payload["station_profile_id"] = cp.Station_profile_id
	_payload["type"] = "adif"
	_payload["string"] = cp.Payload

	// Convert payload to JSON
	_jsonPayload, err := json.Marshal(_payload)
	if err != nil {
		agg_logger.Get().Log("Error marshaling JSON:", err.Error())
		return false
	}

	// Create the HTTP request
	_req, err := http.NewRequest("POST", "https://n7akg.cloudlog.co.uk/index.php/api/qso", bytes.NewBuffer(_jsonPayload))
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
