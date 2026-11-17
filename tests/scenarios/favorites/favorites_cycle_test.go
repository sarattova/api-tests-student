package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"

	"github.com/stretchr/testify/require"
)

func TestFavoritesLifecycle(t *testing.T) {
	visitor := support.LoginPrimary(t)
	owner := support.LoginSecondary(t)
	_, ad := support.CreateAdvertisementFixture(t, owner.Token)

	t.Log("Шаг: добавляем чужое объявление в избранное")
	added := favoriteManager.AddFavorite(t, visitor.Token, ad.ID)
	require.Equalf(t, http.StatusCreated, added.StatusCode, "ошибка добавления в избранное: %+v", added.Error)
	require.Equal(t, ad.ID, added.Data.Advertisement.ID, "в избранное добавилось другое объявление")
	t.Cleanup(func() {
		cleanup := favoriteManager.DeleteFavorite(t, visitor.Token, ad.ID)
		if cleanup.StatusCode != http.StatusNoContent && cleanup.StatusCode != http.StatusNotFound {
			t.Errorf("очистка: не удалось удалить объявление %s из избранного, статус %d, ошибка %+v",
				ad.ID, cleanup.StatusCode, cleanup.Error)
		}
	})

	t.Log("Проверка: объявление появилось в GET /favorites")
	listed := favoriteManager.GetFavorites(t, visitor.Token, ad.Title)
	require.Equalf(t, http.StatusOK, listed.StatusCode, "ошибка получения списка избранного: %+v", listed.Error)
	require.Contains(t, assertions.FavoriteAdvertisementIDs(listed.Data.Items), ad.ID,
		"добавленное объявление отсутствует в избранном")

	t.Log("Шаг: удаляем объявление из избранного")
	removed := favoriteManager.DeleteFavorite(t, visitor.Token, ad.ID)
	require.Equalf(t, http.StatusNoContent, removed.StatusCode, "ошибка удаления из избранного: %+v", removed.Error)
	listed = favoriteManager.GetFavorites(t, visitor.Token, ad.Title)
	require.Equal(t, http.StatusOK, listed.StatusCode, "не удалось получить избранное после удаления записи")
	require.NotContains(t, assertions.FavoriteAdvertisementIDs(listed.Data.Items), ad.ID,
		"удалённое объявление осталось в избранном")

	t.Log("Шаг: добавляем снова, затем владелец удаляет объявление")
	added = favoriteManager.AddFavorite(t, visitor.Token, ad.ID)
	require.Equalf(t, http.StatusCreated, added.StatusCode, "ошибка повторного добавления в избранное: %+v", added.Error)
	deleted := advertisements.DeleteAdvertisement(t, owner.Token, ad.ID)
	require.Equalf(t, http.StatusNoContent, deleted.StatusCode, "ошибка удаления объявления: %+v", deleted.Error)

	t.Log("Проверка: после удаления объявления оно исчезло из избранного")
	listed = favoriteManager.GetFavorites(t, visitor.Token, ad.Title)
	require.Equal(t, http.StatusOK, listed.StatusCode, "не удалось получить избранное после удаления объявления")
	require.NotContains(t, assertions.FavoriteAdvertisementIDs(listed.Data.Items), ad.ID,
		"удалённое объявление осталось в избранном")
	t.Log("Результат: добавление, удаление из избранного и удаление объявления проверены")
}
