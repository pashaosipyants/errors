package errors

type setup struct {
	skip int
	annotation *string
	extender func(error) error
}

func (s *setup) setup(opts []Option) *setup {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Option func(*setup)

func Annotation(annotation string) Option {
	return func(o *setup) {
		o.annotation = &annotation
	}
}

func Extender(extender func(error) error) Option {
	return func(o *setup) {
		o.extender = extender
	}
}

func Skip(skip int) Option {
	return func(o *setup) {
		o.skip = skip
	}
}

func SkipAdd(add int) Option {
	return func(o *setup) {
		o.skip += add
	}
}
