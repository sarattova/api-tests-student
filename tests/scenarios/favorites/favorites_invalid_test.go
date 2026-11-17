package favorites

import (
	"net/http"
	"testing"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"
)

func TestFavoritesRejectEmptyAdvertisementID(t *testing.T) {
	user := support.LoginPrimary(t)

	t.Log("Шаг: пытаемся добавить в избранное объявление без ID")
	result := favoriteManager.AddFavorite(t, user.Token, "")
	t.Logf("Результат: HTTP %d, сообщение %q", result.StatusCode, result.Error.Message)
	support.RequireAPIError(t, result, http.StatusBadRequest, "invalid UUID length")
}
