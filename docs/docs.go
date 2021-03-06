// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/operations/deposit": {
            "post": {
                "description": "increases user's balance by given id and money amount, and returns operation info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "operations"
                ],
                "summary": "increases user's balance",
                "parameters": [
                    {
                        "description": "Operation parameters (receiver id is redundant)",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.OperationInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.Operation"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    }
                }
            }
        },
        "/operations/transfer": {
            "post": {
                "description": "decreases initiator user's balance and increases receiver's balance, and returns operation info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "operations"
                ],
                "summary": "transfers money from one user to another",
                "parameters": [
                    {
                        "description": "Operation parameters",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.OperationInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.Operation"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    }
                }
            }
        },
        "/operations/withdraw": {
            "post": {
                "description": "decreases user's balance by given id and money amount, and returns operation info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "operations"
                ],
                "summary": "decreases user's balance",
                "parameters": [
                    {
                        "description": "Operation parameters (receiver id is redundant)",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.OperationInput"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Withdraw currency",
                        "name": "currency",
                        "in": "query"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.Operation"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    }
                }
            }
        },
        "/users/balance": {
            "post": {
                "description": "returns user's money amount by given id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "shows user's balance",
                "parameters": [
                    {
                        "description": "User ID (amount is redundant)",
                        "name": "id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    }
                }
            }
        },
        "/users/history": {
            "post": {
                "description": "returns a list of operations in which the user appeared, starting from the end",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "returns user's history of operations",
                "parameters": [
                    {
                        "description": "History input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.HistoryInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.RepositoryOperation"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorJSON"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.ErrorJSON": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "domain.HistoryInput": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "mode": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "domain.Operation": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "initiator": {
                    "$ref": "#/definitions/domain.User"
                },
                "receiver": {
                    "$ref": "#/definitions/domain.User"
                },
                "timestamp": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "domain.OperationInput": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "initiator_id": {
                    "type": "integer"
                },
                "receiver_id": {
                    "type": "integer"
                }
            }
        },
        "domain.RepositoryOperation": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "initiator_id": {
                    "type": "integer"
                },
                "receiver_id": {
                    "type": "integer"
                },
                "timestamp": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "domain.User": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8000",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Balance control API",
	Description: "This is a multi-user balance control system.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
