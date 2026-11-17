package advertisementCreation

import (
	"net/http"
	"strings"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/support"
)

func TestCreateAdvertisementRejectsInvalidInput(t *testing.T) {
	user := support.LoginPrimary(t)
	photo := support.PhotoPath(t)
	valid := adModels.CreateAdvertisementRequest{
		Title:       support.UniqueTitle(t),
		Description: "Проверка некорректных параметров",
		Price:       100,
		Quantity:    1,
	}

	cases := []struct {
		name        string
		change      func(*adModels.CreateAdvertisementRequest)
		photos      []string
		messagePart string
	}{
		{"пустое название", func(ad *adModels.CreateAdvertisementRequest) { ad.Title = "" }, []string{photo}, "title is required"},
		{"пустое описание", func(ad *adModels.CreateAdvertisementRequest) { ad.Description = "" }, []string{photo}, "description is required"},
		{"отрицательная цена", func(ad *adModels.CreateAdvertisementRequest) { ad.Price = -1 }, []string{photo}, "price must be between"},
		{"отрицательное количество", func(ad *adModels.CreateAdvertisementRequest) { ad.Quantity = -1 }, []string{photo}, "quantity must be between"},
		{"слишком длинное название", func(ad *adModels.CreateAdvertisementRequest) { ad.Title = strings.Repeat("x", 301) }, []string{photo}, "title must not exceed 50"},
		{"без фотографии", func(*adModels.CreateAdvertisementRequest) {}, nil, "at least 1 photo is required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := valid
			tc.change(&request)
			t.Logf("Шаг: создаём объявление с ошибкой: %s", tc.name)
			result := advertisements.CreateAdvertisement(t, user.Token, request, tc.photos)
			if result.StatusCode == http.StatusCreated && result.Data.ID != "" {
				t.Cleanup(func() {
					cleanup := advertisements.DeleteAdvertisement(t, user.Token, result.Data.ID)
					if cleanup.StatusCode != http.StatusNoContent {
						t.Errorf("очистка: неожиданно созданное объявление %s не удалено, статус %d", result.Data.ID, cleanup.StatusCode)
					}
				})
			}
			support.RequireAPIError(t, result, http.StatusBadRequest, tc.messagePart)
			t.Logf("Результат: запрос отклонён, HTTP %d, сообщение %q", result.StatusCode, result.Error.Message)
		})
	}
}
