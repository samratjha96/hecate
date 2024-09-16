package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

func truncateString(s string, maxLen int) string {
	asRunes := []rune(s)

	if len(asRunes) > maxLen {
		return string(asRunes[:maxLen])
	}

	return s
}

func decodeJsonFromRequest[T any](client RequestDoer, request *http.Request) (T, error) {
	response, err := client.Do(request)
	var result T

	if err != nil {
		return result, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return result, err
	}

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf(
			"unexpected status code %d for %s, response: %s",
			response.StatusCode,
			request.URL,
			truncateString(string(body), 256),
		)
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}
