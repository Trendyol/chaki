package response

type Entity[T any] struct {
	data    T
	status  int
	headers map[string]string
}

func NewEntity[T any](status int, data T, headers map[string]string) Entity[T] {
	return Entity[T]{
		status:  status,
		data:    data,
		headers: headers,
	}
}

func (e Entity[T]) Status() int {
	return e.status
}

func (e Entity[T]) Headers() map[string]string {
	return e.headers
}

func (e Entity[T]) ToResponse() any {
	return Success(e.data)
}
