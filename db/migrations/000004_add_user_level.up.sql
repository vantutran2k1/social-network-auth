CREATE TYPE USER_LEVEL AS ENUM (
    'BRONZE',
    'SILVER',
    'GOLD'
);

ALTER TABLE users
ADD COLUMN level USER_LEVEL NOT NULL DEFAULT 'BRONZE';