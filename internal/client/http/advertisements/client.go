package advertisements

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).Status(http.StatusOK).End().Response
}

// HttpPostAdvertisement отправляет запрос со всеми полями объявления и фотографиями.
func HttpPostAdvertisement(t *testing.T, token string, ad models.CreateAdvertisementRequest, photoPaths []string) *http.Response {
	t.Helper()
	request := apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", ad.Title).
		MultipartFormData("description", ad.Description).
		MultipartFormData("price", strconv.FormatFloat(ad.Price, 'f', -1, 64)).
		MultipartFormData("quantity", strconv.Itoa(ad.Quantity))
	if len(photoPaths) > 0 {
		request = request.MultipartFile("photos", photoPaths...)
	}
	return request.Expect(t).End().Response
}

// HttpGetAdvertisement получает объявление по параметру id в строке запроса API.
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление авторизованного владельца.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления.
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	photoPath := fmt.Sprintf(path.AdvertisementPhotosPath, url.PathEscape(id))
	return apiRunner.GetRunner().Auth(token).Create().Get(photoPath).
		Expect(t).End().Response
}
