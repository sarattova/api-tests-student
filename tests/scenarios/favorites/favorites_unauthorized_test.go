package favorites

import (
	"net/http"
	"testing"

	favoriteManager "api-tests-template/internal/managers/favorites"
	base "api-tests-template/tests"
	"api-tests-template/tests/support"
)

func TestFavoritesRequireAuthorization(t *testing.T) {
	base.SetupSuite()
	for _, tc := range []struct {
		name        string
		token       string
		messagePart string
	}{
		{"без токена", "", "authorization header required"},
		{"неверный токен", "invalid_token", "invalid or expired token"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("список", func(t *testing.T) {
				t.Log("Шаг: GET /favorites без действующей авторизации")
				result := favoriteManager.GetFavorites(t, tc.token, "")
				support.RequireAPIError(t, result, http.StatusUnauthorized, tc.messagePart)
			})
			t.Run("добавление", func(t *testing.T) {
				t.Log("Шаг: POST /favorites без действующей авторизации")
				result := favoriteManager.AddFavorite(t, tc.token, missingAdvertisementID)
				support.RequireAPIError(t, result, http.StatusUnauthorized, tc.messagePart)
			})
			t.Run("удаление", func(t *testing.T) {
				t.Log("Шаг: DELETE /favorites/{ad_id} без действующей авторизации")
				result := favoriteManager.DeleteFavorite(t, tc.token, missingAdvertisementID)
				support.RequireAPIError(t, result, http.StatusUnauthorized, tc.messagePart)
			})
		})
	}
}
