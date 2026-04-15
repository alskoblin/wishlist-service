package domain

import "time"

type Item struct {
	ID          int64
	WishlistID  int64
	Title       string
	Description string
	ProductURL  string
	Priority    int
	Reserved    bool
	ReservedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
