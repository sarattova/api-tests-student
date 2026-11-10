package assertions

import (
	"testing"

	"api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	require.Failf(t, "Объявление не найдено", "Название %q отсутствует в результатах поиска", title)
}
