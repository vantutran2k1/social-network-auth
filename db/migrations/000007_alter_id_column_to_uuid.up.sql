CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Step 1: Drop foreign constraints to users table
ALTER TABLE profiles
DROP CONSTRAINT profiles_user_id_fkey;

ALTER TABLE tokens
DROP CONSTRAINT tokens_user_id_fkey;

ALTER TABLE password_reset_tokens
DROP CONSTRAINT password_reset_tokens_user_id_fkey;

-- Step 1: Create new id for tables
    -- users
ALTER TABLE users
ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();

UPDATE users
SET new_id = uuid_generate_v4()
WHERE new_id IS NULL;

ALTER TABLE users
DROP CONSTRAINT users_pkey;

ALTER TABLE users
ADD CONSTRAINT users_pkey PRIMARY KEY (new_id);

    -- profiles
ALTER TABLE profiles
ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();

UPDATE profiles
SET new_id = uuid_generate_v4()
WHERE new_id IS NULL;

ALTER TABLE profiles
DROP CONSTRAINT profiles_pkey;

ALTER TABLE profiles
ADD CONSTRAINT profiles_pkey PRIMARY KEY (new_id);

ALTER TABLE profiles
DROP COLUMN id;

ALTER TABLE profiles
RENAME COLUMN new_id to id;

    -- tokens
ALTER TABLE tokens
ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();

UPDATE tokens
SET new_id = uuid_generate_v4()
WHERE new_id IS NULL;

ALTER TABLE tokens
DROP CONSTRAINT tokens_pkey;

ALTER TABLE tokens
ADD CONSTRAINT tokens_pkey PRIMARY KEY (new_id);

ALTER TABLE tokens
DROP COLUMN id;

ALTER TABLE tokens
RENAME COLUMN new_id to id;

    -- password_reset_tokens
ALTER TABLE password_reset_tokens
ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();

UPDATE password_reset_tokens
SET new_id = uuid_generate_v4()
WHERE new_id IS NULL;

ALTER TABLE password_reset_tokens
DROP CONSTRAINT password_reset_tokens_pkey;

ALTER TABLE password_reset_tokens
ADD CONSTRAINT password_reset_tokens_pkey PRIMARY KEY (new_id);

ALTER TABLE password_reset_tokens
DROP COLUMN id;

ALTER TABLE password_reset_tokens
RENAME COLUMN new_id to id;

-- Step 3: Update foreign keys of all foreign tables
    -- profiles
ALTER TABLE profiles
ADD COLUMN new_user_id UUID;

UPDATE profiles
SET new_user_id = users.new_id
FROM users
WHERE profiles.user_id = users.id;

ALTER TABLE profiles
ADD CONSTRAINT profiles_user_id_fkey
FOREIGN KEY (new_user_id) REFERENCES users(new_id);

ALTER TABLE profiles
DROP COLUMN user_id;

ALTER TABLE profiles
RENAME COLUMN new_user_id TO user_id;

    -- tokens
ALTER TABLE tokens
ADD COLUMN new_user_id UUID;

UPDATE tokens
SET new_user_id = users.new_id
FROM users
WHERE tokens.user_id = users.id;

ALTER TABLE tokens
ADD CONSTRAINT tokens_user_id_fkey
FOREIGN KEY (new_user_id) REFERENCES users(new_id);

ALTER TABLE tokens
DROP COLUMN user_id;

ALTER TABLE tokens
RENAME COLUMN new_user_id TO user_id;

    -- password_reset_tokens
ALTER TABLE password_reset_tokens
ADD COLUMN new_user_id UUID;

UPDATE password_reset_tokens
SET new_user_id = users.new_id
FROM users
WHERE password_reset_tokens.user_id = users.id;

ALTER TABLE password_reset_tokens
ADD CONSTRAINT password_reset_tokens_user_id_fkey
FOREIGN KEY (new_user_id) REFERENCES users(new_id);

ALTER TABLE password_reset_tokens
DROP COLUMN user_id;

ALTER TABLE password_reset_tokens
RENAME COLUMN new_user_id TO user_id;

-- Step 4: Update id column of users table
ALTER TABLE users
DROP COLUMN id;

ALTER TABLE users
RENAME COLUMN new_id TO id;