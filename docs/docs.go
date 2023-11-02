// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Dmitry Kolesnikov (Zalimannard)",
            "email": "zalimannard@mail.ru"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/albums": {
            "get": {
                "description": "Retrieves a list of all albums, including their best covers if requested.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Retrieve all albums",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of best covers for each album to retrieve",
                        "name": "bestCovers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response with a list of albums and optional best covers for each",
                        "schema": {
                            "$ref": "#/definitions/album_handler.getAllResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid bestCovers format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/albums/{albumId}": {
            "get": {
                "description": "Retrieves detailed information about an album, including its best covers if requested.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Retrieve album details",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Album ID",
                        "name": "albumId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Number of best covers to retrieve",
                        "name": "bestCovers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/album_handler.getResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid albumId or bestCovers format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "404": {
                        "description": "Album not found",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/albums/{albumId}/songs": {
            "get": {
                "description": "Retrieves all songs that are part of the specified album, including detailed information about each song.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Songs"
                ],
                "summary": "Retrieve songs by album ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Unique identifier of the album",
                        "name": "albumId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with a list of songs belonging to the requested album",
                        "schema": {
                            "$ref": "#/definitions/song_handler.getByAlbumIdResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid albumId format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "404": {
                        "description": "Album not found",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/artists": {
            "get": {
                "description": "Retrieves a list of all artists, including their best covers if requested.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Artists"
                ],
                "summary": "Retrieve all artists",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of best covers for each artist to retrieve",
                        "name": "bestCovers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response with a list of artists and optional best covers for each",
                        "schema": {
                            "$ref": "#/definitions/artist_handler.getAllResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid bestCovers format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/artists/{artistId}": {
            "get": {
                "description": "Retrieves detailed information about an artist, including its best covers if requested.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Artists"
                ],
                "summary": "Retrieve artist details",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Artist ID",
                        "name": "artistId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Number of best covers to retrieve",
                        "name": "bestCovers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/artist_handler.getResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid artistId or bestCovers format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "404": {
                        "description": "Artist not found",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/artists/{artistId}/songs": {
            "get": {
                "description": "Retrieves all songs that are part of the specified artist, including detailed information about each song.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Songs"
                ],
                "summary": "Retrieve songs by artist ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Unique identifier of the artist",
                        "name": "artistId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with a list of songs belonging to the requested artist",
                        "schema": {
                            "$ref": "#/definitions/song_handler.getByArtistIdResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid artistId format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "404": {
                        "description": "Artist not found",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/genres": {
            "get": {
                "description": "Retrieves a list of all genres, including their best covers if requested.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Genres"
                ],
                "summary": "Retrieve all genres",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of best covers for each genre to retrieve",
                        "name": "bestCovers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response with a list of genres and optional best covers for each",
                        "schema": {
                            "$ref": "#/definitions/genre_handler.getAllResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid bestCovers format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/genres/{genreId}": {
            "get": {
                "description": "Retrieves detailed information about a genre, including its best covers if requested.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Genres"
                ],
                "summary": "Retrieve genre details",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Genre ID",
                        "name": "genreId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Number of best covers to retrieve",
                        "name": "bestCovers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/genre_handler.getResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid genreId or bestCovers format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "404": {
                        "description": "Genre not found",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/genres/{genreId}/songs": {
            "get": {
                "description": "Retrieves all songs that are part of the specified genre, including detailed information about each song.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Songs"
                ],
                "summary": "Retrieve songs by genre ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Unique identifier of the genre",
                        "name": "genreId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with a list of songs belonging to the requested genre",
                        "schema": {
                            "$ref": "#/definitions/song_handler.getByGenreIdResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid genreId format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "404": {
                        "description": "Genre not found",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/scan": {
            "post": {
                "description": "Scans the system for any new or updated songs, updating the database accordingly.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Scan"
                ],
                "summary": "Initiate a scan for new or updated songs",
                "responses": {
                    "200": {
                        "description": "Successfully initiated song scan"
                    },
                    "500": {
                        "description": "Internal Server Error due to failure in scanning process",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/songs": {
            "get": {
                "description": "Retrieves detailed information about all available songs.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Songs"
                ],
                "summary": "Retrieve a list of all songs",
                "responses": {
                    "200": {
                        "description": "Successful response with list of songs",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/song_handler.getAllResponseItem"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        },
        "/songs/{songId}": {
            "get": {
                "description": "Retrieves detailed information about a song specified by its unique ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Songs"
                ],
                "summary": "Retrieve a song by its ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Unique identifier of the song",
                        "name": "songId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with song details",
                        "schema": {
                            "$ref": "#/definitions/song_handler.getResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid songId format",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "404": {
                        "description": "Song not found",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "album_handler.getAllResponse": {
            "type": "object",
            "properties": {
                "albums": {
                    "description": "Array of albums.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/album_handler.getAllResponseItem"
                    }
                }
            }
        },
        "album_handler.getAllResponseItem": {
            "type": "object",
            "properties": {
                "albumId": {
                    "description": "Unique identifier for the album.",
                    "type": "integer"
                },
                "bestCovers": {
                    "description": "Optional array of best cover IDs for the album.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "title": {
                    "description": "Title of the album.",
                    "type": "string"
                }
            }
        },
        "album_handler.getResponse": {
            "type": "object",
            "properties": {
                "albumId": {
                    "description": "Unique identifier for the album.",
                    "type": "integer"
                },
                "bestCovers": {
                    "description": "Optional array of best cover IDs.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "title": {
                    "description": "Title of the album.",
                    "type": "string"
                }
            }
        },
        "artist_handler.getAllResponse": {
            "type": "object",
            "properties": {
                "artists": {
                    "description": "Array of artists.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/artist_handler.getAllResponseItem"
                    }
                }
            }
        },
        "artist_handler.getAllResponseItem": {
            "type": "object",
            "properties": {
                "artistId": {
                    "description": "Unique identifier for the artist.",
                    "type": "integer"
                },
                "bestCovers": {
                    "description": "Optional array of best cover IDs for the artist.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "name": {
                    "description": "Name of the artist.",
                    "type": "string"
                }
            }
        },
        "artist_handler.getResponse": {
            "type": "object",
            "properties": {
                "artistId": {
                    "description": "Unique identifier for the artist.",
                    "type": "integer"
                },
                "bestCovers": {
                    "description": "Optional array of best cover IDs.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "name": {
                    "description": "Name of the artist.",
                    "type": "string"
                }
            }
        },
        "genre_handler.getAllResponse": {
            "type": "object",
            "properties": {
                "genres": {
                    "description": "Array of genres.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/genre_handler.getAllResponseItem"
                    }
                }
            }
        },
        "genre_handler.getAllResponseItem": {
            "type": "object",
            "properties": {
                "bestCovers": {
                    "description": "Optional array of best cover IDs for the genre.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "genreId": {
                    "description": "Unique identifier for the genre.",
                    "type": "integer"
                },
                "name": {
                    "description": "Name of the genre.",
                    "type": "string"
                }
            }
        },
        "genre_handler.getResponse": {
            "type": "object",
            "properties": {
                "bestCovers": {
                    "description": "Optional array of best cover IDs.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "genreId": {
                    "description": "Unique identifier for the genre.",
                    "type": "integer"
                },
                "name": {
                    "description": "Name of the genre.",
                    "type": "string"
                }
            }
        },
        "response.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "description": "Human-readable error message",
                    "type": "string"
                },
                "reason": {
                    "description": "Internal error description",
                    "type": "string"
                }
            }
        },
        "song_handler.getAllResponseItem": {
            "type": "object",
            "properties": {
                "albumId": {
                    "description": "AlbumId is the identifier of the album to which the song belongs.",
                    "type": "integer"
                },
                "artistId": {
                    "description": "ArtistId is the identifier of the song's artist.",
                    "type": "integer"
                },
                "audioFileId": {
                    "description": "AudioFileId is the identifier of the associated audio file.",
                    "type": "integer"
                },
                "discNumber": {
                    "description": "DiscNumber is the disc number of the song in the album.",
                    "type": "integer"
                },
                "genreId": {
                    "description": "GenreId is the genre identifier of the song.",
                    "type": "integer"
                },
                "lyrics": {
                    "description": "Lyrics are the lyrics of the song.",
                    "type": "string"
                },
                "sha256": {
                    "description": "Sha256 is the SHA256 hash of the song file.",
                    "type": "string"
                },
                "songId": {
                    "description": "SongId is the unique identifier for the song.",
                    "type": "integer"
                },
                "songNumber": {
                    "description": "SongNumber is the track number of the song in the album.",
                    "type": "integer"
                },
                "title": {
                    "description": "Title is the title of the song.",
                    "type": "string"
                },
                "year": {
                    "description": "Year is the release year of the song.",
                    "type": "integer"
                }
            }
        },
        "song_handler.getByAlbumIdResponse": {
            "type": "object",
            "properties": {
                "songs": {
                    "description": "Array of songs belonging to a specific album.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/song_handler.getByAlbumIdResponseItem"
                    }
                }
            }
        },
        "song_handler.getByAlbumIdResponseItem": {
            "type": "object",
            "properties": {
                "albumId": {
                    "description": "Identifier of the album to which the song belongs.",
                    "type": "integer"
                },
                "artistId": {
                    "description": "Identifier of the artist of the song.",
                    "type": "integer"
                },
                "audioFileId": {
                    "description": "Identifier for the associated audio file.",
                    "type": "integer"
                },
                "discNumber": {
                    "description": "Disc number of the song in the album.",
                    "type": "integer"
                },
                "genreId": {
                    "description": "Genre identifier of the song.",
                    "type": "integer"
                },
                "lyrics": {
                    "description": "Lyrics of the song.",
                    "type": "string"
                },
                "sha256": {
                    "description": "SHA256 hash of the song file.",
                    "type": "string"
                },
                "songId": {
                    "description": "Unique identifier for the song.",
                    "type": "integer"
                },
                "songNumber": {
                    "description": "Track number of the song in the album.",
                    "type": "integer"
                },
                "title": {
                    "description": "Title of the song.",
                    "type": "string"
                },
                "year": {
                    "description": "Release year of the song.",
                    "type": "integer"
                }
            }
        },
        "song_handler.getByArtistIdResponse": {
            "type": "object",
            "properties": {
                "songs": {
                    "description": "Array of songs belonging to a specific artist.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/song_handler.getByArtistIdResponseItem"
                    }
                }
            }
        },
        "song_handler.getByArtistIdResponseItem": {
            "type": "object",
            "properties": {
                "albumId": {
                    "description": "Identifier of the album to which the song belongs.",
                    "type": "integer"
                },
                "artistId": {
                    "description": "Identifier of the artist of the song.",
                    "type": "integer"
                },
                "audioFileId": {
                    "description": "Identifier for the associated audio file.",
                    "type": "integer"
                },
                "discNumber": {
                    "description": "Disc number of the song in the album.",
                    "type": "integer"
                },
                "genreId": {
                    "description": "Genre identifier of the song.",
                    "type": "integer"
                },
                "lyrics": {
                    "description": "Lyrics of the song.",
                    "type": "string"
                },
                "sha256": {
                    "description": "SHA256 hash of the song file.",
                    "type": "string"
                },
                "songId": {
                    "description": "Unique identifier for the song.",
                    "type": "integer"
                },
                "songNumber": {
                    "description": "Track number of the song in the album.",
                    "type": "integer"
                },
                "title": {
                    "description": "Title of the song.",
                    "type": "string"
                },
                "year": {
                    "description": "Release year of the song.",
                    "type": "integer"
                }
            }
        },
        "song_handler.getByGenreIdResponse": {
            "type": "object",
            "properties": {
                "songs": {
                    "description": "Array of songs belonging to a specific artist.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/song_handler.getByGenreIdResponseItem"
                    }
                }
            }
        },
        "song_handler.getByGenreIdResponseItem": {
            "type": "object",
            "properties": {
                "albumId": {
                    "description": "Identifier of the album to which the song belongs.",
                    "type": "integer"
                },
                "artistId": {
                    "description": "Identifier of the artist of the song.",
                    "type": "integer"
                },
                "audioFileId": {
                    "description": "Identifier for the associated audio file.",
                    "type": "integer"
                },
                "discNumber": {
                    "description": "Disc number of the song in the album.",
                    "type": "integer"
                },
                "genreId": {
                    "description": "Genre identifier of the song.",
                    "type": "integer"
                },
                "lyrics": {
                    "description": "Lyrics of the song.",
                    "type": "string"
                },
                "sha256": {
                    "description": "SHA256 hash of the song file.",
                    "type": "string"
                },
                "songId": {
                    "description": "Unique identifier for the song.",
                    "type": "integer"
                },
                "songNumber": {
                    "description": "Track number of the song in the album.",
                    "type": "integer"
                },
                "title": {
                    "description": "Title of the song.",
                    "type": "string"
                },
                "year": {
                    "description": "Release year of the song.",
                    "type": "integer"
                }
            }
        },
        "song_handler.getResponse": {
            "type": "object",
            "properties": {
                "albumId": {
                    "description": "AlbumId is the identifier of the album to which the song belongs.",
                    "type": "integer"
                },
                "artistId": {
                    "description": "ArtistId is the identifier of the song's artist.",
                    "type": "integer"
                },
                "audioFileId": {
                    "description": "AudioFileId is the identifier of the associated audio file.",
                    "type": "integer"
                },
                "discNumber": {
                    "description": "DiscNumber is the disc number of the song in the album.",
                    "type": "integer"
                },
                "genreId": {
                    "description": "GenreId is the genre identifier of the song.",
                    "type": "integer"
                },
                "lyrics": {
                    "description": "Lyrics are the lyrics of the song.",
                    "type": "string"
                },
                "sha256": {
                    "description": "Sha256 is the SHA256 hash of the song file.",
                    "type": "string"
                },
                "songId": {
                    "description": "SongId is the unique identifier for the song.",
                    "type": "integer"
                },
                "songNumber": {
                    "description": "SongNumber is the track number of the song in the album.",
                    "type": "integer"
                },
                "title": {
                    "description": "Title is the title of the song.",
                    "type": "string"
                },
                "year": {
                    "description": "Year is the release year of the song.",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.4",
	Host:             "localhost:8023",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Wakarimi Music Metadata API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
