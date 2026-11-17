package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	client "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	shared "api-tests-template/internal/managers/models"

	"github.com/stretchr/testify/require"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := client.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// CreateAdvertisement отправляет запрос на создание и разбирает ответ с объявлением или ошибкой API.
func CreateAdvertisement(t *testing.T, token string, ad models.CreateAdvertisementRequest, photoPaths []string) shared.Result[models.Advertisement] {
	t.Helper()
	return shared.DecodeResponse[models.Advertisement](t, client.HttpPostAdvertisement(t, token, ad, photoPaths))
}

// GetAdvertisement возвращает объявление по идентификатору и HTTP-статус ответа.
func GetAdvertisement(t *testing.T, id string) shared.Result[models.Advertisement] {
	t.Helper()
	return shared.DecodeResponse[models.Advertisement](t, client.HttpGetAdvertisement(t, id))
}

// DeleteAdvertisement удаляет объявление владельца и возвращает статус ответа API.
func DeleteAdvertisement(t *testing.T, token, id string) shared.Result[struct{}] {
	t.Helper()
	return shared.DecodeResponse[struct{}](t, client.HttpDeleteAdvertisement(t, token, id))
}

// GetAdvertisementPhotos возвращает фотографии объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string) shared.Result[[]models.Photo] {
	t.Helper()
	return shared.DecodeResponse[[]models.Photo](t, client.HttpGetAdvertisementPhotos(t, token, id))
}
