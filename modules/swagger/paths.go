package swagger

import (
	"reflect"
	"strings"

	"github.com/Trendyol/chaki/util/slc"
)

func buildPaths(eds []EndpointDef) m {
	p := make(m)
	for _, ed := range eds {
		if strings.HasPrefix(ed.Endpoint, "/__") {
			continue
		}

		desc := m{
			"tags":        []string{ed.Group},
			"summary":     ed.Name,
			"description": "",
			"consumes": []string{
				"application/json",
			},
			"produces": []string{
				"application/json",
			},
			"externalDocs": m{},
		}

		if pr := getParameters(ed.RequestType); pr != nil {
			desc["parameters"] = pr
		}

		desc["responses"] = m{
			"200": m{
				"description": "successful operation",
				"schema":      getPropertyField(ed.ResponseType),
			},
		}

		meth := strings.ToLower(ed.Method)
		swagp := toSwaggerPath(ed.Endpoint)

		existdesc, ok := p[swagp].(m)
		if !ok {
			existdesc = m{}
		}

		existdesc[meth] = desc
		p[swagp] = existdesc

	}

	return p
}

func getParameters(t reflect.Type) []m {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		mi := arrayProperty(t)
		mi["name"] = "body"
		mi["in"] = "body"
		return []m{mi}
	}

	var p []m
	var hasBody bool
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		required := false
		if vts, ok := f.Tag.Lookup("validate"); ok {
			if slc.Contains(strings.Split(vts, ","), "required") {
				required = true
			}
		}

		if n := f.Tag.Get("param"); n != "" {
			pi := getPropertyField(f.Type)
			pi["in"] = "path"
			pi["name"] = n
			pi["description"] = ""
			pi["required"] = true
			p = append(p, pi)
		}

		if n := f.Tag.Get("query"); n != "" {
			pi := getPropertyField(f.Type)
			pi["in"] = "query"
			pi["name"] = n
			pi["description"] = ""
			if required {
				pi["required"] = true
			}
			p = append(p, pi)
		}

		if f.Tag.Get("json") != "" {
			hasBody = true
		}

	}

	if hasBody {
		p = append(p, m{
			"in":          "body",
			"name":        "body",
			"description": "",
			"required":    true,
			"schema": m{
				refKey: withDefinitionPrefix(getNameFromType(t)),
			},
		})
	}

	return p
}
