package favorites

import (
	"net/http"
	"testing"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"
)

const missingAdvertisementID = "00000000-0000-4000-8000-000000000000"

func TestFavoritesRejectAddingNonexistentAdvertisement(t *testing.T) {
	user := support.LoginPrimary(t)

	t.Log("Шаг: добавляем несуществующий ID")
	result := favoriteManager.AddFavorite(t, user.Token, missingAdvertisementID)
	support.RequireAPIError(t, result, http.StatusNotFound, "advertisement not found")
	t.Logf("Результат: HTTP %d, %q", result.StatusCode, result.Error.Message)
}
