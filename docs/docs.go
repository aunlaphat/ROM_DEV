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
        "/orders/allgetorder": {
            "get": {
                "description": "Get all Order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Get Order",
                "operationId": "Allget-order",
                "responses": {
                    "200": {
                        "description": "Order Get",
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
                                                "$ref": "#/definitions/entity.Order"
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
        "/orders/create-order": {
            "post": {
                "description": "Create a new order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Create Order",
                "operationId": "create-order",
                "parameters": [
                    {
                        "description": "Order Data",
                        "name": "order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateOrderRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Order Created",
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
                                                "$ref": "#/definitions/entity.Order"
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
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/orders/delete/{orderNo}": {
            "delete": {
                "description": "Delete an order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Delete Order",
                "operationId": "delete-order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order Number",
                        "name": "orderNo",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Order Deleted",
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
                    "204": {
                        "description": "No Content, Order Delete Successfully",
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
                    "404": {
                        "description": "Order Not Found",
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
        "/orders/getbyID/{orderNo}": {
            "get": {
                "description": "Get details of an order by its order number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Get Order by ID",
                "operationId": "get-order-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order Number",
                        "name": "orderNo",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Order Get by ID",
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
                                                "$ref": "#/definitions/entity.Order"
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
                        "description": "Order not found",
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
        "/orders/update/{orderNo}": {
            "put": {
                "description": "Update an existing order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Update Order",
                "operationId": "update-order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order Number",
                        "name": "orderNo",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated Order Data",
                        "name": "order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateOrderRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Order Updated",
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
                    "204": {
                        "description": "No Content, Order Updated Successfully",
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
                    "404": {
                        "description": "Order Not Found",
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
        "entity.Order": {
            "type": "object",
            "properties": {
                "brandName": {
                    "type": "string",
                    "example": "BEWELL"
                },
                "createDate": {
                    "type": "string",
                    "example": "20/11/2567"
                },
                "custAddress": {
                    "type": "string",
                    "example": "7/20 ซอย15/1"
                },
                "custDistrict": {
                    "type": "string",
                    "example": "บางกรวย"
                },
                "custName": {
                    "type": "string",
                    "example": "Fa"
                },
                "custPhoneNum": {
                    "type": "string",
                    "example": "0921234567"
                },
                "custPostCode": {
                    "type": "string",
                    "example": "11130"
                },
                "custProvince": {
                    "type": "string",
                    "example": "นนทบุรี"
                },
                "custSubDistrict": {
                    "type": "string",
                    "example": "บางกรวย"
                },
                "orderLine": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.OrderLine"
                    }
                },
                "orderNo": {
                    "type": "string",
                    "example": "AB0001"
                },
                "updateDate": {
                    "type": "string",
                    "example": "20/11/2568"
                },
                "userCreated": {
                    "type": "string",
                    "example": "intern"
                },
                "userUpdates": {
                    "type": "string",
                    "example": "intern"
                }
            }
        },
        "entity.OrderLine": {
            "type": "object",
            "properties": {
                "itemName": {
                    "type": "string",
                    "example": "เก้าอี้"
                },
                "orderNo": {
                    "type": "string",
                    "example": "AB0001"
                },
                "price": {
                    "type": "number",
                    "example": 199.05
                },
                "qty": {
                    "type": "integer"
                },
                "sku": {
                    "type": "string"
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
        "request.CreateOrderRequest": {
            "type": "object",
            "properties": {
                "brandName": {
                    "type": "string",
                    "example": "BEWELL"
                },
                "custAddress": {
                    "type": "string",
                    "example": "7/20 ซอย15/1"
                },
                "custDistrict": {
                    "type": "string",
                    "example": "บางกรวย"
                },
                "custName": {
                    "type": "string",
                    "example": "Fa"
                },
                "custPhoneNum": {
                    "type": "string",
                    "example": "0921234567"
                },
                "custPostCode": {
                    "type": "string",
                    "example": "11130"
                },
                "custProvince": {
                    "type": "string",
                    "example": "นนทบุรี"
                },
                "custSubDistrict": {
                    "type": "string",
                    "example": "บางกรวย"
                },
                "orderLines": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/request.OrderLineRequest"
                    }
                },
                "orderNo": {
                    "type": "string",
                    "example": "AB0001"
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
        "request.OrderLineRequest": {
            "type": "object",
            "properties": {
                "itemName": {
                    "type": "string",
                    "example": "เก้าอี้"
                },
                "price": {
                    "type": "number",
                    "example": 199.05
                },
                "qty": {
                    "type": "integer"
                },
                "sku": {
                    "type": "string"
                }
            }
        },
        "request.UpdateOrderRequest": {
            "type": "object",
            "properties": {
                "custAddress": {
                    "type": "string",
                    "example": "7/20 ซอย15/1"
                },
                "custDistrict": {
                    "type": "string",
                    "example": "บางกรวย"
                },
                "custName": {
                    "type": "string",
                    "example": "Fa"
                },
                "custPhoneNum": {
                    "type": "string",
                    "example": "0921234567"
                },
                "custPostCode": {
                    "type": "string",
                    "example": "11130"
                },
                "custProvince": {
                    "type": "string",
                    "example": "นนทบุรี"
                },
                "custSubDistrict": {
                    "type": "string",
                    "example": "บางกรวย"
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
