package main

import (
	"io"
	"net/http"
)

type RequestConfig struct {
	Method  string
	URL     string
	Body    io.Reader
	Headers map[string]string
}

func SendRequest(config RequestConfig) (string, error) {
	if config.Method == "" {
		config.Method = "GET"
	}
	request, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		return "", err
	}
	for key, value := range config.Headers {
		request.Header.Set(key, value)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
