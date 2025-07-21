package ad

import "time"

type Ad struct {
	ID          int     	`json:"id"`
	Title       string  	`json:"title"`
	Description string  	`json:"description"`
	ImageURL    string  	`json:"image_url"`
	Price       float64 	`json:"price"`
	UserID      int     	`json:"user_id"`
	CreatedAt	time.Time 	`json:"created_at"`
}

type CreateAdRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

type FeedAd struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	AuthorLogin string    `json:"author_login"`
	CreatedAt   time.Time `json:"created_at"`
	IsOwner     bool      `json:"is_owner,omitempty"`
}