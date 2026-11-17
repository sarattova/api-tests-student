package models

type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	UserID      string  `json:"user_id"`
	Photos      []Photo `json:"photos"`
}

type AdvertisementsList struct {
	Items      []Advertisement `json:"items"`
	NextCursor string          `json:"next_cursor"`
	HasNext    bool            `json:"has_next"`
}

type Photo struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
}
