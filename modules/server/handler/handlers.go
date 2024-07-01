package handler

import "net/http"

func Get[Req, Res any](path string, h HandlerFunc[Req, Res]) Handler {
	return New[Req, Res](http.MethodGet, path, h)
}

func Post[Req, Res any](path string, h HandlerFunc[Req, Res]) Handler {
	return New[Req, Res](http.MethodPost, path, h)
}

func Patch[Req, Res any](path string, h HandlerFunc[Req, Res]) Handler {
	return New[Req, Res](http.MethodPatch, path, h)
}

func Put[Req, Res any](path string, h HandlerFunc[Req, Res]) Handler {
	return New[Req, Res](http.MethodPut, path, h)
}

func Delete[Req, Res any](path string, h HandlerFunc[Req, Res]) Handler {
	return New[Req, Res](http.MethodDelete, path, h)
}
