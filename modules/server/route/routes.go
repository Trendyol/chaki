package route

import "net/http"

func Get[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New[Req, Res](http.MethodGet, path, h, ds...)
}

func Post[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New[Req, Res](http.MethodPost, path, h, ds...)
}

func Patch[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New[Req, Res](http.MethodPatch, path, h, ds...)
}

func Put[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New[Req, Res](http.MethodPut, path, h, ds...)
}

func Delete[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New[Req, Res](http.MethodDelete, path, h, ds...)
}
