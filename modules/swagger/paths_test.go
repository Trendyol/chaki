package swagger

import (
	"reflect"
	"testing"
)

type request struct {
	ID   string `param:"id"`
	Name string `query:"name"`
	Data string `json:"data"`
}

type anotherRequest struct {
	ID   string   `param:"id"`
	Name []string `query:"names" validate:"required"`
	Data string   `json:"data"`
}

type response struct {
	Field string `json:"field"`
}

func Test_buildPaths(t *testing.T) {
	type args struct {
		eds []EndpointDef
	}
	tests := []struct {
		name string
		args args
		want m
	}{
		{
			name: "Single endpoint with params, query, and json body",
			args: args{
				eds: []EndpointDef{
					{
						Endpoint:     "/test/{id}",
						Group:        "testGroup",
						Name:         "testName",
						Method:       "POST",
						RequestType:  reflect.TypeOf(request{}),
						ResponseType: reflect.TypeOf(response{}),
					},
				},
			},
			want: m{
				"/test/{id}": m{
					"post": m{
						"tags":         []string{"testGroup"},
						"summary":      "testName",
						"description":  "",
						"consumes":     []string{"application/json"},
						"produces":     []string{"application/json"},
						"externalDocs": m{},
						"parameters": []m{
							{
								"in":          "path",
								"name":        "id",
								"description": "",
								"required":    true,
								"type":        "string",
							},
							{
								"in":          "query",
								"name":        "name",
								"description": "",
								"type":        "string",
							},
							{
								"in":          "body",
								"name":        "body",
								"description": "",
								"required":    true,
								"schema": m{
									"$ref": "#/definitions/request",
								},
							},
						},
						"responses": m{
							"200": m{
								"description": "successful operation",
								"schema": m{
									"$ref": "#/definitions/response",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Single endpoint with params, query, and json body with arrays and required fields",
			args: args{
				eds: []EndpointDef{
					{
						Endpoint:     "/test/{id}",
						Group:        "testGroup",
						Name:         "testName",
						Method:       "POST",
						RequestType:  reflect.TypeOf(&anotherRequest{}),
						ResponseType: reflect.TypeOf(&response{}),
					},
				},
			},
			want: m{
				"/test/{id}": m{
					"post": m{
						"tags":         []string{"testGroup"},
						"summary":      "testName",
						"description":  "",
						"consumes":     []string{"application/json"},
						"produces":     []string{"application/json"},
						"externalDocs": m{},
						"parameters": []m{
							{
								"in":          "path",
								"name":        "id",
								"description": "",
								"required":    true,
								"type":        "string",
							},
							{
								"in":          "query",
								"name":        "names",
								"description": "",
								"required":    true,
								"type":        "array",
								"items": m{
									"type": "string",
								},
							},
							{
								"in":          "body",
								"name":        "body",
								"description": "",
								"required":    true,
								"schema": m{
									"$ref": "#/definitions/anotherRequest",
								},
							},
						},
						"responses": m{
							"200": m{
								"description": "successful operation",
								"schema": m{
									"$ref": "#/definitions/response",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildPaths(tt.args.eds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildPaths() = %v\nasdfasdfasddfasfasdfasdf want %v", got, tt.want)
			}
		})
	}
}

func Test_getParameters(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want []m
	}{
		{
			name: "Struct with path and query params",
			args: args{
				t: reflect.TypeOf(request{}),
			},
			want: []m{
				{
					"in":          "path",
					"name":        "id",
					"description": "",
					"required":    true,
					"type":        "string",
				},
				{
					"in":          "query",
					"name":        "name",
					"description": "",
					"type":        "string",
				},
				{
					"in":          "body",
					"name":        "body",
					"description": "",
					"required":    true,
					"schema": m{
						"$ref": "#/definitions/request",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getParameters(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getParameters() = %v\n, want %v", got, tt.want)
			}
		})
	}
}
