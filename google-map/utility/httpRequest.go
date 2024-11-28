package utility

import (
	"bytes"
	"context"
	"google-map/model/enum"
	"net/http"
	"time"
)

func RequestToAnotherService(
	ctx context.Context,
	method enum.RequestMethod,
	url enum.URL,
	headers map[enum.RequestHeader]string,
	body string,
) (*http.Response, error) {

	// create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(ctx, method.ToString(), url.ToString(), bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	// Set headers
	allowedHeaders := []string{"Content-Type", "Accept", "Authorization", "User-Agent"}
	for key, value := range headers {
		req.Header.Set(key.ToString(), value)
	}
	for _, header := range allowedHeaders {
		if val, ok := headers[enum.RequestHeader(header)]; ok {
			req.Header.Set(header, val)
		}
	}

	// Create client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
