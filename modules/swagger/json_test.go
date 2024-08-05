package swagger

import (
	"reflect"
	"testing"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/internal/typlect"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type sampleStruct struct {
	Field string
}

func TestDocs_WithHost(t *testing.T) {
	d := Docs{}
	h := "localhost"
	want := Docs{"Host": h}

	got := d.WithHost(h)
	assert.Equal(t, want, got)
}

func Test_buildDocs(t *testing.T) {
	v := viper.New()
	cfg := config.NewConfig(v, nil)
	eds := []EndpointDef{}

	got := buildDocs(eds, cfg)
	assert.NotNil(t, got)
	assert.Contains(t, got, "swagger")
	assert.Contains(t, got, "definitions")
	assert.Contains(t, got, "paths")
}

func Test_toSwaggerPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/api/:id", "/api/{id}"},
		{"/api/:user/:id", "/api/{user}/{id}"},
	}

	for _, tt := range tests {
		got := toSwaggerPath(tt.input)
		assert.Equal(t, tt.expected, got)
	}
}

func Test_baseJson(t *testing.T) {
	v := viper.New()
	cfg := config.NewConfig(v, nil)
	got := baseJson(cfg)

	assert.NotNil(t, got)
	assert.Equal(t, "2.0", got["swagger"])
	assert.Contains(t, got, "info")
	assert.Contains(t, got, "host")
	assert.Contains(t, got, "basePath")
	assert.Contains(t, got, "schemes")
}

func Test_withDefinitionPrefix(t *testing.T) {
	input := "Example"
	expected := "#/definitions/Example"

	got := withDefinitionPrefix(input)
	assert.Equal(t, expected, got)
}

func Test_getPrimitiveType(t *testing.T) {
	tests := []struct {
		input    reflect.Type
		expected m
	}{
		{reflect.TypeOf(1), m{"type": "integer", "format": "int"}},
		{reflect.TypeOf(true), m{"type": "boolean"}},
		{reflect.TypeOf("string"), m{"type": "string"}},
		{reflect.TypeOf([]int{}), m{"type": "slice"}},
	}

	for _, tt := range tests {
		got := getPrimitiveType(tt.input)
		assert.Equal(t, tt.expected, got)
	}
}

func Test_getPropertyField(t *testing.T) {
	tests := []struct {
		input    reflect.Type
		expected m
	}{
		{typlect.TypeNoParam, m{"type": "string"}},
		{reflect.TypeOf(1), m{"type": "integer", "format": "int"}},
		{reflect.TypeOf(true), m{"type": "boolean"}},
		{reflect.TypeOf([]int{}), m{"type": "array", "items": m{"type": "integer", "format": "int"}}},
		{reflect.TypeOf(sampleStruct{}), m{"$ref": "#/definitions/sampleStruct"}},
	}

	for _, tt := range tests {
		got := getPropertyField(tt.input)
		assert.Equal(t, tt.expected, got)
	}
}

func Test_arrayProperty(t *testing.T) {
	input := reflect.TypeOf([]int{})
	expected := m{"type": "array", "items": m{"type": "integer", "format": "int"}}

	got := arrayProperty(input)
	assert.Equal(t, expected, got)
}

func Test_getNameFromType(t *testing.T) {
	input := reflect.TypeOf(sampleStruct{})
	expected := "sampleStruct"

	got := getNameFromType(input)
	assert.Equal(t, expected, got)
}
