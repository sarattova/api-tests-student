package advertisementCreation

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/support"

	"github.com/stretchr/testify/require"
)

func TestCreateAdvertisementWithAllFields(t *testing.T) {
	user := support.LoginPrimary(t)
	request, created := support.CreateAdvertisementFixture(t, user.Token)

	t.Log("Проверка: ответ на создание содержит все переданные поля")
	require.Equal(t, request.Title, created.Title, "после создания название не совпадает")
	require.Equal(t, request.Description, created.Description, "после создания описание не совпадает")
	require.Equal(t, request.Price, created.Price, "после создания цена не совпадает")
	require.Equal(t, request.Quantity, created.Quantity, "после создания количество не совпадает")
	require.Equal(t, user.User.Id, created.UserID, "после создания ID владельца не совпадает")
	require.Len(t, created.Photos, 1, "у созданного объявления должна быть одна загруженная фотография")
	require.NotEmpty(t, created.Photos[0].ID, "у загруженной фотографии отсутствует ID")
	require.NotEmpty(t, created.Photos[0].URL, "у загруженной фотографии отсутствует URL")

	t.Log("Шаг: получаем созданное объявление по ID")
	byID := advertisements.GetAdvertisement(t, created.ID)
	require.Equalf(t, http.StatusOK, byID.StatusCode, "ошибка получения объявления: %+v", byID.Error)
	require.Equal(t, created.ID, byID.Data.ID, "при получении по ID вернулся другой ID объявления")
	require.Equal(t, request.Title, byID.Data.Title, "при получении по ID название не совпадает")
	require.Equal(t, request.Description, byID.Data.Description, "при получении по ID описание не совпадает")
	require.Equal(t, request.Price, byID.Data.Price, "при получении по ID цена не совпадает")
	require.Equal(t, request.Quantity, byID.Data.Quantity, "при получении по ID количество не совпадает")
	require.Equal(t, user.User.Id, byID.Data.UserID, "при получении по ID владелец не совпадает")

	t.Log("Шаг: ищем объявление по полному названию")
	search := advertisements.GetAdvertisementsList(t, user.Token, request.Title)
	found := false
	for _, item := range search.Items {
		if item.ID == created.ID {
			found = true
			require.Equal(t, request.Title, item.Title, "в поиске найдено объявление с другим названием")
			break
		}
	}
	require.Truef(t, found, "объявление %s с названием %q не найдено в результатах поиска", created.ID, request.Title)

	t.Log("Шаг: получаем загруженные фотографии")
	photos := advertisements.GetAdvertisementPhotos(t, user.Token, created.ID)
	require.Equalf(t, http.StatusOK, photos.StatusCode, "ошибка получения фотографий: %+v", photos.Error)
	require.Len(t, photos.Data, 1, "ожидалась ровно одна загруженная фотография")
	require.Equal(t, created.Photos[0].ID, photos.Data[0].ID, "ID фотографии из отдельной ручки не совпадает")
	require.Equal(t, created.Photos[0].URL, photos.Data[0].URL, "URL фотографии из отдельной ручки не совпадает")
	t.Log("Результат: поля объявления, результат поиска и фотография совпадают")
}
