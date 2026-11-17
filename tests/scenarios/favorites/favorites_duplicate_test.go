package favorites

import (
	"net/http"
	"testing"

	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"

	"github.com/stretchr/testify/require"
)

func TestFavoritesRejectDuplicate(t *testing.T) {
	visitor := support.LoginPrimary(t)
	owner := support.LoginSecondary(t)
	_, ad := support.CreateAdvertisementFixture(t, owner.Token)

	t.Log("Подготовка: добавляем чужое объявление в избранное")
	first := favoriteManager.AddFavorite(t, visitor.Token, ad.ID)
	require.Equalf(t, http.StatusCreated, first.StatusCode, "ошибка первого добавления: %+v", first.Error)
	t.Cleanup(func() {
		cleanup := favoriteManager.DeleteFavorite(t, visitor.Token, ad.ID)
		if cleanup.StatusCode != http.StatusNoContent && cleanup.StatusCode != http.StatusNotFound {
			t.Errorf("очистка: не удалось удалить объявление %s из избранного, статус %d", ad.ID, cleanup.StatusCode)
		}
	})

	t.Log("Шаг: повторно добавляем то же объявление")
	duplicate := favoriteManager.AddFavorite(t, visitor.Token, ad.ID)
	support.RequireAPIError(t, duplicate, http.StatusConflict, "already in favorites")
	t.Logf("Результат: повторное добавление отклонено с сообщением %q", duplicate.Error.Message)
}
