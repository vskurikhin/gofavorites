// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "fiber@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/login": {
            "post": {
                "security": [
                    {
                        "none": []
                    }
                ],
                "description": "аутентификация пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "аутентификация",
                "parameters": [
                    {
                        "description": "Формат запроса JSON (body)",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.SignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "пользователь успешно аутентифицирован",
                        "schema": {
                            "$ref": "#/definitions/dto.SignInRequest"
                        }
                    },
                    "400": {
                        "description": "неверный формат запроса",
                        "schema": {
                            "$ref": "#/definitions/dto.SignInRequest"
                        }
                    },
                    "401": {
                        "description": "неверная пара логин/пароль",
                        "schema": {
                            "$ref": "#/definitions/dto.SignInRequest"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/favorites/get": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    },
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "избранное получения инструментов для пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Favorites"
                ],
                "summary": "избранное",
                "responses": {
                    "200": {
                        "description": "успешная обработка запроса",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/dto.Favorites"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "пользователь не авторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "none": []
                    }
                ],
                "description": "избранное получения инструмента для пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Favorites"
                ],
                "summary": "избранное",
                "parameters": [
                    {
                        "description": "Формат запроса JSON (body)",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Favorites"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "получение инструмента",
                        "schema": {
                            "$ref": "#/definitions/dto.Favorites"
                        }
                    },
                    "400": {
                        "description": "неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "пользователь не авторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/favorites/set": {
            "post": {
                "security": [
                    {
                        "none": []
                    }
                ],
                "description": "избранное сохранение инструмента для пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Favorites"
                ],
                "summary": "избранное",
                "parameters": [
                    {
                        "description": "Формат запроса JSON (body)",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Favorites"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "получение инструмента",
                        "schema": {
                            "$ref": "#/definitions/dto.Favorites"
                        }
                    },
                    "400": {
                        "description": "неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "пользователь не авторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.Favorites": {
            "type": "object",
            "required": [
                "asset_type",
                "isin"
            ],
            "properties": {
                "asset_type": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "isin": {
                    "type": "string"
                }
            }
        },
        "dto.SignInRequest": {
            "type": "object",
            "required": [
                "password",
                "user_name"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "security": [
        {
            "Bearer": []
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8443",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "GoFavorites API",
	Description:      "This is a sample swagger for Fiber",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
