CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wishlists (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    event_date TIMESTAMPTZ NOT NULL,
    public_token TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wishlist_items (
    id BIGSERIAL PRIMARY KEY,
    wishlist_id BIGINT NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    product_url TEXT NOT NULL DEFAULT '',
    priority INTEGER NOT NULL DEFAULT 1,
    reserved BOOLEAN NOT NULL DEFAULT FALSE,
    reserved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wishlists_owner_id ON wishlists(owner_id);
CREATE INDEX IF NOT EXISTS idx_wishlist_items_wishlist_id ON wishlist_items(wishlist_id);
