-- +goose Up
 CREATE TABLE profiles (
     id UUID PRIMARY KEY,
     name VARCHAR(255) UNIQUE NOT NULL,
     gender VARCHAR(50) NOT NULL,
     gender_probability FLOAT NOT NULL,
     sample_size INTEGER NOT NULL,
     age INTEGER NOT NULL,
     age_group VARCHAR(50) NOT NULL,
     country_id VARCHAR(10) NOT NULL,
     country_probability FLOAT NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMPTZ
 );

 CREATE INDEX idx_profiles_gender ON profiles (gender);
 CREATE INDEX idx_profiles_country_id ON profiles (country_id);
 CREATE INDEX idx_profiles_age_group ON profiles (age_group);

-- +goose Down
 DROP INDEX IF EXISTS idx_profiles_gender;
 DROP INDEX IF EXISTS idx_profiles_country_id;
 DROP INDEX IF EXISTS idx_profiles_age_group;

 DROP TABLE IF EXISTS profiles;
