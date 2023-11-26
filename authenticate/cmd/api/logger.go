package main

import (
	"bytes"
	"net/http"
)

func (s *Server) LogRequest(name, data string) error {
	request, err := http.NewRequest(
		"POST",
		"http://logger-service/write_log",
		bytes.NewBuffer([]byte(`{"name":"`+name+`","data":"`+data+`"}`)))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
