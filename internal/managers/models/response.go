package models

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Result[T any] struct {
	StatusCode int
	Data       T
	Error      ErrorResponse
}

// DecodeResponse разбирает ответ API как успешную модель или модель ошибки.
func DecodeResponse[T any](t *testing.T, response *http.Response) Result[T] {
	t.Helper()
	require.NotNil(t, response, "ответ API отсутствует")
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать ответ API")

	result := Result[T]{StatusCode: response.StatusCode}
	if len(body) == 0 {
		return result
	}
	if response.StatusCode >= http.StatusBadRequest {
		require.NoError(t, json.Unmarshal(body, &result.Error), "не удалось разобрать ошибку API")
		return result
	}
	require.NoError(t, json.Unmarshal(body, &result.Data), "не удалось разобрать ответ API")
	return result
}
