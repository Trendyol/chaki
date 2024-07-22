package consumer

type Starter interface {
	Start() error
	Stop() error
}
