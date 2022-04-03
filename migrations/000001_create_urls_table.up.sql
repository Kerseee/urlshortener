CREATE TABLE IF NOT EXISTS urls (
    id bigserial PRIMARY KEY,
    url text  NOT NULL,
    short_url text UNIQUE NOT NULL,
    expire_at timestamp with time zone NOT NULL
);