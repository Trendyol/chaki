package swagger

import (
	"reflect"
)

type EndpointDef struct {
	RequestType  reflect.Type
	ResponseType reflect.Type
	Group        string
	Name         string
	Endpoint     string
	Method       string
}
