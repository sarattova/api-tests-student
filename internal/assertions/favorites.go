package assertions

import "api-tests-template/internal/managers/favorites/models"

// FavoriteAdvertisementIDs возвращает идентификаторы объявлений для проверок в сценариях.
func FavoriteAdvertisementIDs(items []models.FavoriteItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.Advertisement.ID)
	}
	return ids
}
