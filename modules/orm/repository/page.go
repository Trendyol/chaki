package repository

type PageableResponse[T any] struct {
	Page       int
	Size       int
	Data       []T
	TotalCount int
	TotalPage  int
}

type PageableRequest struct {
	Page int
	Size int
	Sort []string
}
