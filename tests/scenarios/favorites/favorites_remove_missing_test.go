package favorites

import (
	"net/http"
	"testing"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"
)

func TestFavoritesRejectRemovingNonexistentAdvertisement(t *testing.T) {
	user := support.LoginPrimary(t)

	t.Log("Шаг: удаляем несуществующий ID")
	result := favoriteManager.DeleteFavorite(t, user.Token, missingAdvertisementID)
	support.RequireAPIError(t, result, http.StatusNotFound, "favorite entry not found")
	t.Logf("Результат: HTTP %d, %q", result.StatusCode, result.Error.Message)
}
