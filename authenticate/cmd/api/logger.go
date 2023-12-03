package main

import (
	"authenticate/rabbitmq/event"
	"bytes"
	"encoding/json"
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

// rabbitmq

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (s *Server) logEventViaRabbit(name, data, qName string) error {
	emitter, err := event.NewEventEmitter(s.rabbitMQ)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: data,
	}

	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), qName) //"log.INFO", "log.ERROR", "log.CRITICAL"
	if err != nil {
		return err
	}
	return nil
}
