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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/injector/ClearSidecar": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Remove the sidecar from a given deployment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "injector"
                ],
                "summary": "Remove the sidecar",
                "parameters": [
                    {
                        "description": "ClearSidecarPayload type",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/injectormodels.ClearSidecarPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/injectormodels.Deployment"
                        }
                    }
                }
            }
        },
        "/api/injector/GetDeployments": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get deployments for a given namespace or all namespaces",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "injector"
                ],
                "summary": "Obtain a list of Deployment objects",
                "parameters": [
                    {
                        "description": "GetDeploymentsPayload type",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/injectormodels.GetDeploymentsPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/injectormodels.Deployment"
                            }
                        }
                    }
                }
            }
        },
        "/api/injector/GetSingleDeployment": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get deployments for a given namespace and name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "injector"
                ],
                "summary": "Obtain a specific Deployment objects",
                "parameters": [
                    {
                        "description": "GetSingleDeploymentPayload type",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/injectormodels.GetSingleDeploymentPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/injectormodels.Deployment"
                        }
                    }
                }
            }
        },
        "/api/injector/SetSidecar": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Set the sidecar for a given deployment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "injector"
                ],
                "summary": "Activate the sidecar",
                "parameters": [
                    {
                        "description": "SetSidecarPayload type",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/injectormodels.SetSidecarPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/injectormodels.Deployment"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "injectormodels.ClearSidecarPayload": {
            "type": "object",
            "properties": {
                "DeploymentName": {
                    "type": "string"
                },
                "Namespace": {
                    "type": "string"
                },
                "SidecarContainerName": {
                    "type": "string"
                }
            }
        },
        "injectormodels.Deployment": {
            "type": "object",
            "properties": {
                "Name": {
                    "type": "string"
                },
                "Namespace": {
                    "type": "string"
                },
                "VolumeNames": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "injectormodels.GetDeploymentsPayload": {
            "type": "object",
            "properties": {
                "Filtered": {
                    "type": "boolean"
                },
                "Namespace": {
                    "type": "string"
                }
            }
        },
        "injectormodels.GetSingleDeploymentPayload": {
            "type": "object",
            "properties": {
                "DeploymentName": {
                    "type": "string"
                },
                "Namespace": {
                    "type": "string"
                }
            }
        },
        "injectormodels.SetSidecarPayload": {
            "type": "object",
            "properties": {
                "Command": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "DeploymentName": {
                    "type": "string"
                },
                "Namespace": {
                    "type": "string"
                },
                "SidecarContainerName": {
                    "type": "string"
                },
                "SidecarImage": {
                    "type": "string"
                },
                "VolumeMounts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/injectormodels.Volume"
                    }
                }
            }
        },
        "injectormodels.Volume": {
            "type": "object",
            "properties": {
                "MountPath": {
                    "type": "string"
                },
                "Name": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "X-API-KEY",
            "in": "header"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Kubernetes OnDemand Sidecar Injector API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
