BEGIN;

CREATE TABLE claims (
    id BIGSERIAL PRIMARY KEY,
    created_by UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL,
    status TEXT NOT NULL,
    photos TEXT[],
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    feedback TEXT DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    status_updated_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    feedback_updated_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE subcategories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    category_id BIGINT REFERENCES categories(id) ON DELETE CASCADE
);

COMMIT;