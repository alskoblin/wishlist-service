package dto

type CreateItemInput struct {
	WishlistID  int64  `json:"-"`
	OwnerID     int64  `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductURL  string `json:"product_url"`
	Priority    int    `json:"priority"`
}

type UpdateItemInput struct {
	ItemID      int64  `json:"-"`
	WishlistID  int64  `json:"-"`
	OwnerID     int64  `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductURL  string `json:"product_url"`
	Priority    int    `json:"priority"`
}

type ReserveItemInput struct {
	PublicToken string `json:"public_token"`
	ItemID      int64  `json:"item_id"`
}
