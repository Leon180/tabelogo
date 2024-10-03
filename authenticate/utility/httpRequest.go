package utility

import (
	"authenticate/model/enum"
	"bytes"
	"net/http"
	"time"
)

func RequestToAnotherService(
	method enum.RequestMethod,
	url enum.URL,
	headers map[enum.RequestHeader]string,
	body string,
) (*http.Response, error) {

	// Create client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// Create request
	req, err := http.NewRequest(method.ToString(), url.ToString(), bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	// Set headers
	for key, value := range headers {
		req.Header.Set(key.ToString(), value)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
