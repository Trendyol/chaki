package swagger

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/internal/typlect"
)

type m map[string]any

type Docs m

func (d Docs) WithHost(h string) Docs {
	d["Host"] = h
	return d
}

const (
	refKey = "$ref"
)

func buildDocs(eds []EndpointDef, cfg *config.Config) Docs {
	dj := baseJSON(cfg)
	dj["definitions"] = buildDefinitions(eds)
	dj["paths"] = buildPaths(eds)
	return Docs(dj)
}

var pathRegexp = regexp.MustCompile(`(\:[A-Za-z0-9_]*)`)

func toSwaggerPath(s string) string {
	return pathRegexp.ReplaceAllStringFunc(s, func(s string) string {
		return fmt.Sprintf("{%s}", s[1:])
	})
}

func baseJSON(cfg *config.Config) m {
	of := cfg.Of("info")

	of.SetDefault("description", "")
	of.SetDefault("version", "1.0.0")
	of.SetDefault("title", "Application")
	of.SetDefault("terms", "http://swagger.io/terms/")
	of.SetDefault("email", "")

	return m{
		"swagger": "2.0",
		"info": m{
			"description":    of.GetString("description"),
			"version":        of.GetString("version"),
			"title":          of.GetString("title"),
			"termsOfService": of.GetString("terms"),
			"contact": m{
				"email": of.GetString("email"),
			},
			"license": m{
				"name": "Apache 2.0",
				"url":  "http://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		"host":     "",
		"basePath": "/",
		"schemes":  []string{},
	}
}

func withDefinitionPrefix(s string) string {
	return fmt.Sprintf("#/definitions/%s", s)
}

func getPrimitiveType(t reflect.Type) m {
	if kp := t.Kind().String(); strings.HasPrefix(kp, "int") {
		return m{
			"type":   "integer",
			"format": kp,
		}
	}

	k := t.Kind().String()

	if t.Kind() == reflect.Bool {
		k = "boolean"
	}

	return m{
		"type": k,
	}
}

func getPropertyField(t reflect.Type) m {
	if t == typlect.TypeNoParam {
		return m{"type": "string"}
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t == typlect.TypeTime {
		b, _ := time.Now().MarshalJSON()
		return m{"type": "string", "example": strings.Trim(string(b), "\"")}
	}

	if t.Kind() == reflect.Struct {
		return m{
			refKey: withDefinitionPrefix(getNameFromType(t)),
		}
	}

	if t.Kind() == reflect.Slice {
		return arrayProperty(t)
	}

	return getPrimitiveType(t)
}

func arrayProperty(t reflect.Type) m {
	it := t.Elem()
	if it.Kind() == reflect.Pointer {
		it = it.Elem()
	}

	return m{
		"type":  "array",
		"items": getPropertyField(it),
	}
}

func getNameFromType(t reflect.Type) string {
	s := strings.ReplaceAll(t.Name(), "]", "")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "*", "")
	return strings.ReplaceAll(s, "[", "__")
}
