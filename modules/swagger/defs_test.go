package swagger

import (
	"reflect"
	"testing"

	"github.com/Trendyol/chaki/internal/typlect"
)

type testStruct struct {
	Field1 string          `json:"field1" validate:"required"`
	Field2 int             `json:"field2"`
	Field3 anotherStruct   `json:"field3"`
	Field4 []anotherStruct `json:"field4"`
	Field5 *string         `json:"field5"`
	Field6 map[string]int  `json:"field6" validate:"required"`
}

type anotherStruct struct {
	Field1 string `json:"field1"`
	Field2 []int  `json:"field2"`
}

func TestBuildDefinitions(t *testing.T) {
	eds := []EndpointDef{
		{
			RequestType:  reflect.TypeOf(testStruct{}),
			ResponseType: reflect.TypeOf(anotherStruct{}),
		},
		{
			RequestType:  reflect.TypeOf(&testStruct{}),
			ResponseType: reflect.TypeOf(&anotherStruct{}),
		},
	}

	defs := buildDefinitions(eds)

	if len(defs) != 2 {
		t.Errorf("Expected 2 definitions, got %d", len(defs))
	}

	if _, ok := defs["testStruct"]; !ok {
		t.Errorf("Expected testStruct in definitions, not found")
	}

	if _, ok := defs["anotherStruct"]; !ok {
		t.Errorf("Expected anotherStruct in definitions, not found")
	}
}

func TestBuildModelDefinition(t *testing.T) {
	mockType := reflect.TypeOf(testStruct{})
	defs := make(m)

	buildModelDefinition(defs, mockType, true)

	if len(defs) == 0 {
		t.Errorf("Expected definitions to be populated, got empty map")
	}

	if _, ok := defs["testStruct"]; !ok {
		t.Errorf("Expected testStruct in definitions, not found")
	}

	requiredFields := defs["testStruct"].(m)["required"].([]string)
	if !reflect.DeepEqual(requiredFields, []string{"field1", "field6"}) {
		t.Errorf("Expected required fields to be populated, got %+v", requiredFields)
	}

	if props, ok := defs["testStruct"].(m)["properties"]; ok {
		if p, ok := props.(m)["field1"]; !ok {
			t.Errorf("Expected field1 in properties, not found")
		} else {
			if fieldType, ok := p.(m)["type"]; !ok || fieldType != "string" {
				t.Errorf("Expected field1 to be of type string")
			}
		}
		if p, ok := props.(m)["field2"]; !ok {
			t.Errorf("Expected field2 in properties, not found")
		} else {
			if fieldType, ok := p.(m)["type"]; !ok || fieldType != "integer" {
				t.Errorf("Expected field2 to be of type integer")
			}
		}

		if p, ok := props.(m)["field3"]; !ok {
			t.Errorf("Expected Field3 in properties, not found")
		} else {
			if fieldType, ok := p.(m)["$ref"]; !ok || fieldType != "#/definitions/anotherStruct" {
				t.Errorf("Expected field3 to be of type ref")
			}
		}

		if p, ok := props.(m)["field4"]; !ok {
			t.Errorf("Expected field4 in properties, not found")
		} else {
			if fieldType, ok := p.(m)["type"]; !ok || fieldType != "array" {
				t.Errorf("Expected field4 to be of type array")
			}
		}

		if p, ok := props.(m)["field5"]; !ok {
			t.Errorf("Expected field5 in properties, not found")
		} else {
			if fieldType, ok := p.(m)["type"]; !ok || fieldType != "string" {
				t.Errorf("Expected field5 to be of type string")
			}
		}

		if p, ok := props.(m)["field6"]; !ok {
			t.Errorf("Expected field6 in properties, not found")
		} else {
			if fieldType, ok := p.(m)["type"]; !ok || fieldType != "map" {
				t.Errorf("Expected field6 to be of type map")
			}
		}
	} else {
		t.Errorf("Expected properties to be a map, got %T", defs["testStruct"].(m)["properties"])
	}
}

func TestBuildModelDefinitions(t *testing.T) {
	tests := []struct {
		field reflect.Type
		len   int
		name  string
	}{
		{typlect.TypeNoParam, 0, "no param"},
		{reflect.TypeOf([]string{}), 0, "slice"},

		// As test struct contains one inner struct, there should be 2 different definitions
		{reflect.TypeOf(&testStruct{}), 2, "pointer"},
		{reflect.TypeOf(testStruct{}), 2, "struct"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defs := make(m)
			buildModelDefinition(defs, tt.field, true)

			if len(defs) != tt.len {
				t.Errorf("Expected %d definitions, got %d", tt.len, len(defs))
			}
		})
	}
}

func TestGetFieldName(t *testing.T) {
	tests := []struct {
		field    reflect.StructField
		expected string
	}{
		{reflect.StructField{Name: "Field1", Tag: `json:"field1"`}, "field1"},
		{reflect.StructField{Name: "Field2", Tag: `json:"field2,omitempty"`}, "field2"},
		{reflect.StructField{Name: "Field3"}, "Field3"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := getFieldName(tt.field); got != tt.expected {
				t.Errorf("getFieldName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
