package beater

import (
	"encoding/json"
	"net/http"

	"github.com/elastic/beats/libbeat/logp"
)

// IP represents IP data
type IP struct {
	IP string `json:"ip"`
}

// RetrieveIP gets the IP from ipify.org
func (ip *IP) RetrieveIP() (string, error) {
	req, err := http.NewRequest("GET", "https://api.ipify.org?format=json", nil)

	if err != nil {
		logp.Err("Cannot create request.")

		return "", err
	}

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		logp.Err("Cannot perform request.")

		return "", err
	}

	defer resp.Body.Close()

	errorDecode := json.NewDecoder(resp.Body).Decode(&ip)

	if errorDecode != nil {
		logp.Err("Cannot decode response")

		return "", errorDecode
	}

	return ip.IP, nil
}
