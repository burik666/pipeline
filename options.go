package pipeline

type (
	// Opts represents the options for a pipeline stage.
	Opts struct {
		name string
	}

	// OptFn represents a function that applies an option to the value.
	OptFn = func(*Opts)
)

// Name returns the name of the pipeline stage.
func (o Opts) Name() string {
	return o.name
}

// buildOpts applies the given options to the value.
func buildOpts[T any](opts ...func(*T)) T {
	var v T

	for _, apply := range opts {
		apply(&v)
	}

	return v
}

// WithName sets the name of the pipeline stage.
func WithName(name string) OptFn {
	return func(opts *Opts) {
		opts.name = name
	}
}

// withOpts returns a function that applies the given options to the value.
func withOpts(opts Opts) OptFn {
	return func(o *Opts) {
		*o = opts
	}
}
