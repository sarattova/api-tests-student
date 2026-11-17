package advertisementCreation

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	base "api-tests-template/tests"
	"api-tests-template/tests/support"
)

func TestCreateAdvertisementRequiresAuthorization(t *testing.T) {
	base.SetupSuite()
	photo := support.PhotoPath(t)
	request := adModels.CreateAdvertisementRequest{
		Title:       support.UniqueTitle(t),
		Description: "Проверка запроса без авторизации",
		Price:       100,
		Quantity:    1,
	}

	for _, tc := range []struct {
		name        string
		token       string
		messagePart string
	}{
		{"без токена", "", "authorization header required"},
		{"неверный токен", "invalid_token", "invalid or expired token"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Шаг: пытаемся создать объявление: %s", tc.name)
			result := advertisements.CreateAdvertisement(t, tc.token, request, []string{photo})
			support.RequireAPIError(t, result, http.StatusUnauthorized, tc.messagePart)
			t.Logf("Результат: HTTP %d, %q", result.StatusCode, result.Error.Message)
		})
	}
}
