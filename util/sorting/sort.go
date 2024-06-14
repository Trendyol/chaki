package sorting

type Dir string

const (
	DirAsc  Dir = "asc"
	DirDesc Dir = "desc"
)

type Sorting struct {
	Field string
	Dir   Dir
}

func Asc(field string) Sorting {
	return Sorting{
		Field: field,
		Dir:   DirAsc,
	}
}

func Desc(field string) Sorting {
	return Sorting{
		Field: field,
		Dir:   DirDesc,
	}
}
