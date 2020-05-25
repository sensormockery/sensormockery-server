BEGIN;

CREATE TABLE IF NOT EXISTS wave_types (
    wave_type VARCHAR (16) PRIMARY KEY
);

INSERT INTO wave_types (wave_type) VALUES ('sawtooth');
INSERT INTO wave_types (wave_type) VALUES ('rectangular');
INSERT INTO wave_types (wave_type) VALUES ('triangular');
INSERT INTO wave_types (wave_type) VALUES ('sine');

CREATE TABLE IF NOT EXISTS streams (
    id SERIAL PRIMARY KEY,
    locked INT DEFAULT 0 NOT NULL,
    wave_type VARCHAR (16) NOT NULL REFERENCES wave_types(wave_type) ON DELETE CASCADE,
    sensor VARCHAR (64) NOT NULL,
    noise_coeff REAL DEFAULT 0,
    start_date TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subscription_types (
    subs_type VARCHAR (32) PRIMARY KEY
);

INSERT INTO subscription_types (subs_type) VALUES ('custom');
INSERT INTO subscription_types (subs_type) VALUES ('kafka');

CREATE TABLE IF NOT EXISTS subscriptions (
    subs_endpoint VARCHAR (512) PRIMARY KEY,
    subs_type VARCHAR (32) NOT NULL REFERENCES subscription_types(subs_type) ON DELETE CASCADE,
    stream INT NOT NULL REFERENCES streams(id)
);

COMMIT;
