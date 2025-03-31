package ycloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func request(url string, method string, body []byte, apiKey string) ([]byte, error) {
	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("failed to send with status code: %s", resp.Status))
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}
