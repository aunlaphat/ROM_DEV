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
                "description": "Validates if the JWT token is valid and retrieves user claims.",
                "tags": [
                    "Authentication"
                ],
                "summary": "Check Authentication",
                "responses": {
                    "200": {
                        "description": "Authenticated user details",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Authenticates user credentials and generates a JWT token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "User Login",
                "parameters": [
                    {
                        "description": "User login credentials",
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
                            "$ref": "#/definitions/response.User"
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
                "description": "Authenticates user credentials from Lark and generates a JWT token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "User Lark Login",
                "parameters": [
                    {
                        "description": "User login from Lark",
                        "name": "login-request",
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
                            "$ref": "#/definitions/response.User"
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
                "description": "Logs the user out by removing the JWT token from the cookie.",
                "tags": [
                    "Authentication"
                ],
                "summary": "User Logout",
                "responses": {
                    "200": {
                        "description": "Logout successful",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/order/create": {
            "post": {
                "description": "Creates a new return order including order head and order lines",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Return Order MKP"
                ],
                "summary": "Create a new return order",
                "operationId": "create-return-order",
                "parameters": [
                    {
                        "description": "Return Order Data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateBeforeReturnOrder"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.BeforeReturnOrderResponse"
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
        "/order/search": {
            "get": {
                "description": "Retrieve the details of an order by its SO number or Order number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Return Order MKP"
                ],
                "summary": "Search order by SO number or Order number",
                "operationId": "search-order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SO number",
                        "name": "soNo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Order number",
                        "name": "orderNo",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.SearchOrderResponse"
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
        "/user/{username}": {
            "get": {
                "description": "Get user credentials by userName",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get User Credentials",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UserName",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User credentials",
                        "schema": {
                            "$ref": "#/definitions/response.UserRole"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
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
                "data": {},
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "gin.H": {
            "type": "object",
            "additionalProperties": {}
        },
        "request.CreateBeforeReturnOrder": {
            "type": "object",
            "required": [
                "logistic",
                "mkpStatus",
                "orderNo",
                "returnDate",
                "soNo",
                "soStatus",
                "trackingNo",
                "warehouseID"
            ],
            "properties": {
                "items": {
                    "description": "CreateBy    string                        ` + "`" + `json:\"createBy\" db:\"CreateBy\" binding:\"required\"` + "`" + `",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/request.CreateBeforeReturnOrderItem"
                    }
                },
                "logistic": {
                    "type": "string"
                },
                "mkpStatus": {
                    "type": "string"
                },
                "orderNo": {
                    "type": "string"
                },
                "returnDate": {
                    "description": "Location    string                        ` + "`" + `json:\"location\" db:\"Location\" binding:\"required\"` + "`" + `",
                    "type": "string"
                },
                "soNo": {
                    "type": "string"
                },
                "soStatus": {
                    "description": "SrNo        *string                       ` + "`" + `json:\"srNo,omitempty\" db:\"SrNo\"` + "`" + `",
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
        "request.CreateBeforeReturnOrderItem": {
            "type": "object",
            "required": [
                "createBy",
                "itemName",
                "price",
                "qty",
                "returnQty",
                "sku"
            ],
            "properties": {
                "alterSKU": {
                    "type": "string"
                },
                "createBy": {
                    "type": "string"
                },
                "itemName": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "qty": {
                    "type": "integer"
                },
                "returnQty": {
                    "type": "integer"
                },
                "sku": {
                    "type": "string"
                },
                "trackingNo": {
                    "type": "string"
                }
            }
        },
        "request.LoginLark": {
            "type": "object",
            "properties": {
                "userID": {
                    "type": "string",
                    "example": "DC65060"
                },
                "userName": {
                    "type": "string",
                    "example": "eknarin.ler"
                }
            }
        },
        "request.LoginWeb": {
            "type": "object",
            "properties": {
                "password": {
                    "description": "change password lastest in 17 January 2025",
                    "type": "string",
                    "example": "EKna1234"
                },
                "userName": {
                    "type": "string",
                    "example": "eknarin.ler"
                }
            }
        },
        "response.BeforeReturnOrderItem": {
            "type": "object",
            "properties": {
                "alterSKU": {
                    "type": "string"
                },
                "createBy": {
                    "type": "string"
                },
                "createDate": {
                    "type": "string"
                },
                "itemName": {
                    "type": "string"
                },
                "orderNo": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "qty": {
                    "type": "integer"
                },
                "returnQty": {
                    "type": "integer"
                },
                "sku": {
                    "type": "string"
                },
                "trackingNo": {
                    "type": "string"
                }
            }
        },
        "response.BeforeReturnOrderResponse": {
            "type": "object",
            "properties": {
                "cancelId": {
                    "type": "integer"
                },
                "channelId": {
                    "type": "integer"
                },
                "confirmBy": {
                    "type": "string"
                },
                "confirmDate": {
                    "type": "string"
                },
                "createBy": {
                    "type": "string"
                },
                "createDate": {
                    "type": "string"
                },
                "customerId": {
                    "type": "string"
                },
                "isCNCreated": {
                    "type": "boolean"
                },
                "isEdited": {
                    "type": "boolean"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.BeforeReturnOrderItem"
                    }
                },
                "logistic": {
                    "type": "string"
                },
                "mkpStatus": {
                    "type": "string"
                },
                "orderNo": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                },
                "returnDate": {
                    "type": "string"
                },
                "soNo": {
                    "type": "string"
                },
                "soStatus": {
                    "type": "string"
                },
                "srNo": {
                    "type": "string"
                },
                "statusConfId": {
                    "type": "integer"
                },
                "statusReturnId": {
                    "type": "integer"
                },
                "trackingNo": {
                    "type": "string"
                },
                "updateBy": {
                    "type": "string"
                },
                "updateDate": {
                    "type": "string"
                },
                "warehouseId": {
                    "type": "integer"
                }
            }
        },
        "response.SearchOrderItem": {
            "type": "object",
            "properties": {
                "itemName": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "qty": {
                    "type": "integer"
                },
                "sku": {
                    "type": "string"
                }
            }
        },
        "response.SearchOrderResponse": {
            "type": "object",
            "properties": {
                "createDate": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.SearchOrderItem"
                    }
                },
                "orderNo": {
                    "type": "string"
                },
                "salesStatus": {
                    "type": "string"
                },
                "soNo": {
                    "type": "string"
                },
                "statusMKP": {
                    "type": "string"
                }
            }
        },
        "response.User": {
            "type": "object",
            "properties": {
                "departmentNo": {
                    "type": "string"
                },
                "fullNameTH": {
                    "type": "string"
                },
                "nickName": {
                    "type": "string"
                },
                "platform": {
                    "type": "string"
                },
                "roleID": {
                    "type": "integer"
                },
                "userID": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                }
            }
        },
        "response.UserRole": {
            "type": "object",
            "properties": {
                "departmentNo": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "fullNameTH": {
                    "type": "string"
                },
                "nickName": {
                    "type": "string"
                },
                "permission": {
                    "type": "string"
                },
                "roleID": {
                    "type": "integer"
                },
                "roleName": {
                    "type": "string"
                },
                "userID": {
                    "type": "string"
                },
                "userName": {
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
	Title:            "Return Order Management Service ⭐",
	Description:      "This is a Return Order Management Service API server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
