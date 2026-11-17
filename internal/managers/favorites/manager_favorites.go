package favorites

import (
	"testing"

	client "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
	shared "api-tests-template/internal/managers/models"
)

// GetFavorites возвращает избранное пользователя или ошибку API.
func GetFavorites(t *testing.T, token, search string) shared.Result[models.FavoriteList] {
	t.Helper()
	return shared.DecodeResponse[models.FavoriteList](t, client.HttpGetFavorites(t, token, search))
}

// AddFavorite добавляет объявление в избранное пользователя.
func AddFavorite(t *testing.T, token, advertisementID string) shared.Result[models.FavoriteItem] {
	t.Helper()
	return shared.DecodeResponse[models.FavoriteItem](t, client.HttpPostFavorite(t, token, advertisementID))
}

// DeleteFavorite удаляет объявление из избранного пользователя.
func DeleteFavorite(t *testing.T, token, advertisementID string) shared.Result[struct{}] {
	t.Helper()
	return shared.DecodeResponse[struct{}](t, client.HttpDeleteFavorite(t, token, advertisementID))
}
