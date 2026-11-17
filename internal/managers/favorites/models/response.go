package models

import "api-tests-template/internal/managers/advertisements/models"

type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

type FavoriteItem struct {
	FavoriteID    string               `json:"favorite_id"`
	Advertisement models.Advertisement `json:"advertisement"`
}

type FavoriteList struct {
	Items      []FavoriteItem `json:"items"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}
