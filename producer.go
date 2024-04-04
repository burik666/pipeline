package pipeline

// NewProducer creates a new producer stage.
// Producer is the stage that only produces data.
func NewProducer[T any](fn func(func(T) (T, error)) error, opts ...OptFn) Stage[T] {
	return NewStage(func(in T, next func(T) (T, error)) (T, error) {
		return in, fn(next)
	}, opts...)
}
