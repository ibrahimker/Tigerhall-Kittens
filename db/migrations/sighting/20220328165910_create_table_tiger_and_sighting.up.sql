BEGIN;
CREATE TABLE IF NOT EXISTS sighting.tiger (
    "id" SERIAL not null primary key,
    "name" varchar(255) not null,
    "date_of_birth" date not null,
    "last_seen_timestamp" timestamp not null,
    "last_seen_latitude" numeric not null,
    "last_seen_longitude" numeric not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp,
    "deleted_at" timestamp
);

CREATE INDEX IF NOT EXISTS idx_tiger_last_seen ON sighting.tiger("last_seen_timestamp");

CREATE TABLE IF NOT EXISTS sighting.sighting (
    "id" SERIAL not null primary key,
    "tiger_id" int not null references sighting.tiger(id),
    "seen_at" timestamp not null,
    "latitude" numeric not null,
    "longitude" numeric not null,
    "image_data" text not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp,
    "deleted_at" timestamp
);

CREATE INDEX IF NOT EXISTS idx_sighting_seen_at ON sighting.sighting("seen_at");

COMMIT;