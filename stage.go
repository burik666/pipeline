package pipeline

type (
	// Stage represents a stage in the pipeline.
	Stage[T any] struct {
		fn     StageFn[T]
		nextFn func(T) (T, error)
		// next *Stage[T]

		opts Opts
	}
	// StageFn is the function signature for a stage.
	StageFn[T any] func(in T, next func(T) (T, error)) (T, error)
)

// NewStage creates a new stage with the given function and options.
func NewStage[T any](fn StageFn[T], opts ...OptFn) Stage[T] {
	return Stage[T]{
		fn: fn,

		opts: buildOpts[Opts](opts...),
	}
}

// NewSimpleStage creates a new stage without next argument.
func NewSimpleStage[T any](fn func(T) (T, error), opts ...OptFn) Stage[T] {
	return NewStage(func(in T, next func(T) (T, error)) (T, error) {
		ret, err := fn(in)
		if err != nil {
			return ret, err
		}

		return next(ret)
	}, opts...)
}

// do runs the stage with the given input.
func (s *Stage[T]) do(in T) (T, error) {
	if s.nextFn == nil { // stopStage
		return in, nil
	}

	return s.fn(in, s.nextFn)
}

// stopStage is a helper function to stop the pipeline.
func stopStage[T any](in T, _ func(T) (T, error)) (T, error) {
	return in, nil
}
