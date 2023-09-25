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
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Get all albums",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/album_handler.readAllResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch all album",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/albums/{albumId}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Get detailed information about an album and its tracks by album id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Album Identifier",
                        "name": "albumId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/album_handler.readResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid albumId format",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch album with details",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/artists": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Artists"
                ],
                "summary": "Get all artists",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/artist_handler.readAllResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch all artists",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/artists/{artistId}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Artists"
                ],
                "summary": "Get detailed information about an artist and his tracks by artist id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Artist Identifier",
                        "name": "artistId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/artist_handler.readResponseArtist"
                        }
                    },
                    "400": {
                        "description": "Invalid artistId format",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch artist with details",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/genres": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Genres"
                ],
                "summary": "Get all genres",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/genre_handler.readAllResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch all genres",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/genres/{genreId}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Genres"
                ],
                "summary": "Get detailed information about an genre and his tracks by genre id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Genre Identifier",
                        "name": "genreId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/genre_handler.readResponseGenre"
                        }
                    },
                    "400": {
                        "description": "Invalid genreId format",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch genre with details",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/track-metadata": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TrackMetadata"
                ],
                "summary": "Get all track metadata",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/track_metadata_handler.readAllResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch all genres",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "album_handler.readAllResponse": {
            "description": "List of all albums",
            "type": "object",
            "properties": {
                "albums": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/album_handler.readAllResponseItem"
                    }
                }
            }
        },
        "album_handler.readAllResponseItem": {
            "description": "Album details",
            "type": "object",
            "properties": {
                "albumId": {
                    "type": "integer"
                },
                "coverId": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "tracksCount": {
                    "type": "integer"
                }
            }
        },
        "album_handler.readResponse": {
            "description": "Response structure containing detailed information about an album and its tracks.",
            "type": "object",
            "properties": {
                "albumId": {
                    "type": "integer"
                },
                "coverId": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "trackMetadataList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/album_handler.readResponseTrack"
                    }
                },
                "tracksCount": {
                    "type": "integer"
                }
            }
        },
        "album_handler.readResponseTrack": {
            "description": "Response structure containing details about a track.",
            "type": "object",
            "properties": {
                "albumId": {
                    "type": "integer"
                },
                "artistId": {
                    "type": "integer"
                },
                "discNumber": {
                    "type": "integer"
                },
                "genreId": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "trackId": {
                    "type": "integer"
                },
                "trackMetadataId": {
                    "type": "integer"
                },
                "trackNumber": {
                    "type": "integer"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "artist_handler.readAllResponse": {
            "description": "List of all artists",
            "type": "object",
            "properties": {
                "artists": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/artist_handler.readAllResponseItem"
                    }
                }
            }
        },
        "artist_handler.readAllResponseItem": {
            "description": "Artist details",
            "type": "object",
            "properties": {
                "artistId": {
                    "type": "integer"
                },
                "mostPopularCoversIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "name": {
                    "type": "string"
                },
                "tracksCount": {
                    "type": "integer"
                }
            }
        },
        "artist_handler.readResponseArtist": {
            "description": "Response structure containing detailed information about an artist and his tracks.",
            "type": "object",
            "properties": {
                "artistId": {
                    "type": "integer"
                },
                "mostPopularCoverIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "name": {
                    "type": "string"
                },
                "trackMetadataList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/artist_handler.readResponseTrack"
                    }
                },
                "tracksCount": {
                    "type": "integer"
                }
            }
        },
        "artist_handler.readResponseTrack": {
            "description": "Response structure containing details about a track.",
            "type": "object",
            "properties": {
                "albumId": {
                    "type": "integer"
                },
                "artistId": {
                    "type": "integer"
                },
                "discNumber": {
                    "type": "integer"
                },
                "genreId": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "trackId": {
                    "type": "integer"
                },
                "trackMetadataId": {
                    "type": "integer"
                },
                "trackNumber": {
                    "type": "integer"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "genre_handler.readAllResponse": {
            "description": "List of all genres",
            "type": "object",
            "properties": {
                "genres": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/genre_handler.readAllResponseItem"
                    }
                }
            }
        },
        "genre_handler.readAllResponseItem": {
            "description": "Genre details",
            "type": "object",
            "properties": {
                "genreId": {
                    "type": "integer"
                },
                "mostPopularCoversIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "name": {
                    "type": "string"
                },
                "tracksCount": {
                    "type": "integer"
                }
            }
        },
        "genre_handler.readResponseGenre": {
            "description": "Response structure containing detailed information about an genre and his tracks.",
            "type": "object",
            "properties": {
                "genreId": {
                    "type": "integer"
                },
                "mostPopularCoverIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "name": {
                    "type": "string"
                },
                "trackMetadataList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/genre_handler.readResponseTrack"
                    }
                },
                "tracksCount": {
                    "type": "integer"
                }
            }
        },
        "genre_handler.readResponseTrack": {
            "description": "Response structure containing details about a track.",
            "type": "object",
            "properties": {
                "albumId": {
                    "type": "integer"
                },
                "artistId": {
                    "type": "integer"
                },
                "discNumber": {
                    "type": "integer"
                },
                "genreId": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "trackId": {
                    "type": "integer"
                },
                "trackMetadataId": {
                    "type": "integer"
                },
                "trackNumber": {
                    "type": "integer"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "track_metadata_handler.readAllResponse": {
            "type": "object",
            "properties": {
                "trackMetadataList": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/track_metadata_handler.readAllResponseItem"
                    }
                }
            }
        },
        "track_metadata_handler.readAllResponseAlbum": {
            "type": "object",
            "properties": {
                "albumId": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "track_metadata_handler.readAllResponseArtist": {
            "type": "object",
            "properties": {
                "artistId": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "track_metadata_handler.readAllResponseGenre": {
            "type": "object",
            "properties": {
                "genreId": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "track_metadata_handler.readAllResponseItem": {
            "type": "object",
            "properties": {
                "album": {
                    "$ref": "#/definitions/track_metadata_handler.readAllResponseAlbum"
                },
                "artist": {
                    "$ref": "#/definitions/track_metadata_handler.readAllResponseArtist"
                },
                "discNumber": {
                    "type": "integer"
                },
                "genre": {
                    "$ref": "#/definitions/track_metadata_handler.readAllResponseGenre"
                },
                "title": {
                    "type": "string"
                },
                "trackId": {
                    "type": "integer"
                },
                "trackMetadataId": {
                    "type": "integer"
                },
                "trackNumber": {
                    "type": "integer"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "types.Error": {
            "description": "Response structure for error messages",
            "type": "object",
            "required": [
                "error"
            ],
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.3",
	Host:             "localhost:8023",
	BasePath:         "/api/music-metadata-service",
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
