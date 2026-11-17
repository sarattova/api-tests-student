package favorites

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpGetFavorites запрашивает избранное пользователя с необязательным фильтром по названию.
func HttpGetFavorites(t *testing.T, token, search string) *http.Response {
	t.Helper()
	request := apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath)
	if search != "" {
		request = request.Query("search", search)
	}
	return request.Expect(t).End().Response
}

// HttpPostFavorite добавляет объявление в избранное авторизованного пользователя.
func HttpPostFavorite(t *testing.T, token, advertisementID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		JSON(models.AddFavoriteRequest{AdvertisementID: advertisementID}).
		Expect(t).End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного авторизованного пользователя.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	t.Helper()
	favoritePath := fmt.Sprintf(path.FavoriteByIDPath, url.PathEscape(advertisementID))
	return apiRunner.GetRunner().Auth(token).Create().Delete(favoritePath).
		Expect(t).End().Response
}
