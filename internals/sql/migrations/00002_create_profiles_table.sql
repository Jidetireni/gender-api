-- +goose Up
CREATE TABLE profiles (
    id          UUID PRIMARY KEY,
    name        VARCHAR(255) UNIQUE NOT NULL,
    gender      VARCHAR(50) NOT NULL,
    gender_probability  FLOAT NOT NULL,
    age         INTEGER NOT NULL,
    age_group   VARCHAR(50) NOT NULL,
    country_id  VARCHAR(2) NOT NULL,
    country_name VARCHAR(255) NOT NULL,
    country_probability FLOAT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_profiles_gender     ON profiles (gender);
CREATE INDEX idx_profiles_country_id ON profiles (country_id);
CREATE INDEX idx_profiles_age_group  ON profiles (age_group);
CREATE INDEX idx_profiles_age        ON profiles (age);
CREATE INDEX idx_profiles_created_at ON profiles (created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_profiles_gender;
DROP INDEX IF EXISTS idx_profiles_country_id;
DROP INDEX IF EXISTS idx_profiles_age_group;
DROP INDEX IF EXISTS idx_profiles_age;
DROP INDEX IF EXISTS idx_profiles_created_at;
DROP TABLE IF EXISTS profiles;
