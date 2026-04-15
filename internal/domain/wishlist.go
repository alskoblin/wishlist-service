package domain

import "time"

type Wishlist struct {
	ID          int64
	OwnerID     int64
	EventTitle  string
	Description string
	EventDate   time.Time
	PublicToken string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
