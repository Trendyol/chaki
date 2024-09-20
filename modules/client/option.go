package client

type options struct {
	errDecoder     ErrDecoder
	driverWrappers []DriverWrapper
}

type Option interface {
	Apply(*options)
}

type withOption func(*options)

func (wf withOption) Apply(opts *options) {
	wf(opts)
}

func WithErrDecoder(ed ErrDecoder) Option {
	return withOption(func(o *options) {
		o.errDecoder = ed
	})
}

func WithDriverWrappers(wrappers ...DriverWrapper) Option {
	return withOption(func(o *options) {
		o.driverWrappers = append(o.driverWrappers, wrappers...)
	})
}
