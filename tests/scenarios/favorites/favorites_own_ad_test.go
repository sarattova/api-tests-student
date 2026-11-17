package favorites

import (
	"net/http"
	"testing"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"
)

func TestFavoritesRejectOwnAdvertisement(t *testing.T) {
	owner := support.LoginPrimary(t)
	_, ad := support.CreateAdvertisementFixture(t, owner.Token)

	t.Log("Шаг: пытаемся добавить своё объявление в избранное")
	result := favoriteManager.AddFavorite(t, owner.Token, ad.ID)
	support.RequireAPIError(t, result, http.StatusForbidden, "cannot add own advertisement")
	t.Logf("Результат: добавление своего объявления отклонено с сообщением %q", result.Error.Message)
}
