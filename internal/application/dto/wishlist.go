package dto

import "time"

type CreateWishlistInput struct {
	OwnerID     int64     `json:"-"`
	EventTitle  string    `json:"event_title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
}

type UpdateWishlistInput struct {
	WishlistID  int64     `json:"-"`
	OwnerID     int64     `json:"-"`
	EventTitle  string    `json:"event_title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
}
