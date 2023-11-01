CREATE TABLE "artists"
(
    "artist_id" SERIAL PRIMARY KEY,
    "name"      TEXT NOT NULL UNIQUE
);

CREATE TABLE "albums"
(
    "album_id" SERIAL PRIMARY KEY,
    "title"    TEXT NOT NULL UNIQUE
);

CREATE TABLE "genres"
(
    "genre_id" SERIAL PRIMARY KEY,
    "name"     TEXT NOT NULL UNIQUE
);

CREATE TABLE "songs"
(
    "song_id"       SERIAL PRIMARY KEY,
    "audio_file_id" INTEGER NOT NULL UNIQUE,
    "title"         TEXT,
    "album_id"      INTEGER,
    "artist_id"     INTEGER,
    "genre_id"      INTEGER,
    "year"          INTEGER,
    "song_number"   INTEGER,
    "disc_number"   INTEGER,
    "lyrics"        TEXT,
    "sha_256"       TEXT    NOT NULL UNIQUE,
    FOREIGN KEY ("artist_id") REFERENCES "artists" ("artist_id"),
    FOREIGN KEY ("album_id") REFERENCES "albums" ("album_id"),
    FOREIGN KEY ("genre_id") REFERENCES "genres" ("genre_id")
);
