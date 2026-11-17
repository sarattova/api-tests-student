package assertions

import (
	"testing"

	"api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

// ContainsAdvertisementByTitle проверяет наличие объявления с точным названием в результатах поиска.
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	require.Failf(t, "Объявление не найдено", "Название %q отсутствует в результатах поиска", title)
}
