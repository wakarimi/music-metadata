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
            "name": "Dmitry Kolesnikov (Zalimannard)"
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
                    "albums"
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
                        "description": "Internal Server Error",
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
	Version:          "0.1.2",
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