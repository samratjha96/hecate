package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	maxTruncateLength = 256
)

// RequestDoer is an interface for objects that can execute HTTP requests
type RequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	asRunes := []rune(s)
	if len(asRunes) > maxLen {
		return string(asRunes[:maxLen])
	}
	return s
}

// decodeJSONFromRequest sends an HTTP request and decodes the JSON response
func decodeJSONFromRequest[T any](ctx context.Context, client RequestDoer, request *http.Request) (T, error) {
	var result T

	request = request.WithContext(ctx)
	response, err := client.Do(request)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf(
			"unexpected status code %d for %s, response: %s",
			response.StatusCode,
			request.URL,
			truncateString(string(body), maxTruncateLength),
		)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}
