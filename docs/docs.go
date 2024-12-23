// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/auth": {
            "get": {
                "description": "A test endpoint to check if the user is authenticated and to demonstrate Swagger documentation.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Check Authentication",
                "operationId": "check-authentication",
                "responses": {
                    "200": {
                        "description": "Authenticated user details",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "object",
                                            "additionalProperties": true
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Handles user login requests and generates a token for the authenticated user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User Login",
                "operationId": "user-login",
                "parameters": [
                    {
                        "description": "User login credentials in JSON format",
                        "name": "login-request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.LoginWeb"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT token",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/auth/login-lark": {
            "post": {
                "description": "Handles user login requests and generates a token for the authenticated user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User Lark Login",
                "operationId": "user-login-lark",
                "parameters": [
                    {
                        "description": "User login from lark credentials from Lark in JSON format",
                        "name": "Login-request-lark",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.LoginLark"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT token",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Logs out the user by deleting the JWT token.",
                "tags": [
                    "Auth"
                ],
                "summary": "User Logout",
                "operationId": "user-logout",
                "responses": {
                    "200": {
                        "description": "Logout successful",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/constants/get-district": {
            "get": {
                "description": "Get all Thai District.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Constants"
                ],
                "summary": "Get Thai District",
                "operationId": "get-district",
                "responses": {
                    "200": {
                        "description": "District",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.District"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "District not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/constants/get-province": {
            "get": {
                "description": "Get all Thai Province.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Constants"
                ],
                "summary": "Get Thai Province",
                "operationId": "get-province",
                "responses": {
                    "200": {
                        "description": "Province",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.Province"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Province not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/constants/get-sub-district": {
            "get": {
                "description": "Get all Thai SubDistrict.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Constants"
                ],
                "summary": "Get Thai SubDistrict",
                "operationId": "get-sub-district",
                "responses": {
                    "200": {
                        "description": "SubDistrict",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.SubDistrict"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "SubDistrict not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/return-order/create-return-order": {
            "post": {
                "description": "Create a new return order with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Return Order"
                ],
                "summary": "Create a new return order",
                "operationId": "create-return-order",
                "parameters": [
                    {
                        "description": "Return order details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.BeforeReturnOrderRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/return-order/list-orders": {
            "get": {
                "description": "Retrieve a list of all return orders with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Return Order"
                ],
                "summary": "List all return orders",
                "operationId": "list-orders",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/return-order/{orderNo}": {
            "get": {
                "description": "Retrieve the details of a specific return order by its order number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Return Order"
                ],
                "summary": "Get return order by order number",
                "operationId": "get-order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order number",
                        "name": "orderNo",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/return-order/{orderNo}/cancel": {
            "post": {
                "description": "Cancel a specific return order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Return Order"
                ],
                "summary": "Cancel return order",
                "operationId": "cancel-order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order number",
                        "name": "orderNo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Cancel details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CancelOrderRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/return-order/{orderNo}/status": {
            "put": {
                "description": "Update the status of a specific return order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Return Order"
                ],
                "summary": "Update return order status",
                "operationId": "update-order-status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order number",
                        "name": "orderNo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Status update details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateStatusRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.Response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "result": {},
                "success": {
                    "type": "boolean"
                }
            }
        },
        "entity.District": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "nameEN": {
                    "type": "string"
                },
                "nameTH": {
                    "type": "string"
                },
                "provinceCode": {
                    "type": "integer"
                }
            }
        },
        "entity.Province": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "nameEN": {
                    "type": "string"
                },
                "nameTH": {
                    "type": "string"
                }
            }
        },
        "entity.SubDistrict": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "districtCode": {
                    "type": "integer"
                },
                "nameEN": {
                    "type": "string"
                },
                "nameTH": {
                    "type": "string"
                },
                "zipCode": {
                    "type": "string"
                }
            }
        },
        "request.BeforeReturnOrderLineRequest": {
            "description": "Line item details for return order",
            "type": "object",
            "required": [
                "price",
                "qty",
                "returnQty",
                "sku",
                "trackingNo"
            ],
            "properties": {
                "price": {
                    "type": "number"
                },
                "qty": {
                    "type": "integer",
                    "minimum": 1
                },
                "returnQty": {
                    "type": "integer",
                    "minimum": 1
                },
                "sku": {
                    "type": "string"
                },
                "trackingNo": {
                    "type": "string"
                }
            }
        },
        "request.BeforeReturnOrderRequest": {
            "type": "object",
            "required": [
                "channelID",
                "createBy",
                "customerID",
                "logistic",
                "orderNo",
                "returnDate",
                "returnLines",
                "returnType",
                "saleOrder",
                "saleReturn",
                "trackingNo",
                "warehouseID"
            ],
            "properties": {
                "channelID": {
                    "type": "integer"
                },
                "createBy": {
                    "type": "string"
                },
                "customerID": {
                    "type": "string"
                },
                "logistic": {
                    "type": "string"
                },
                "orderNo": {
                    "type": "string"
                },
                "returnDate": {
                    "type": "string"
                },
                "returnLines": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/request.BeforeReturnOrderLineRequest"
                    }
                },
                "returnType": {
                    "type": "string"
                },
                "saleOrder": {
                    "type": "string"
                },
                "saleReturn": {
                    "type": "string"
                },
                "trackingNo": {
                    "type": "string"
                },
                "warehouseID": {
                    "type": "integer"
                }
            }
        },
        "request.CancelOrderRequest": {
            "type": "object",
            "required": [
                "cancelBy"
            ],
            "properties": {
                "cancelBy": {
                    "type": "string"
                }
            }
        },
        "request.LoginLark": {
            "type": "object",
            "properties": {
                "userID": {
                    "type": "string",
                    "example": "DC99999"
                },
                "userName": {
                    "type": "string",
                    "example": "eknarin"
                }
            }
        },
        "request.LoginWeb": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "example": "asdfhdskjf"
                },
                "userName": {
                    "type": "string",
                    "example": "eknarin"
                }
            }
        },
        "request.UpdateStatusRequest": {
            "type": "object",
            "required": [
                "orderNo",
                "statusId",
                "updateBy"
            ],
            "properties": {
                "orderNo": {
                    "type": "string"
                },
                "statusId": {
                    "type": "integer"
                },
                "updateBy": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Boilerplate Service",
	Description:      "This is a sample server for Boilerplate project .",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
