package pipeline

//go:generate go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@v1.1.0 --output README.md

// Pipeline is a sequence of stages that process some data.
type Pipeline[T any] struct {
	stages []Stage[T]
}

// New creates a new pipeline with the given stages.
func New[T any](stages ...Stage[T]) Pipeline[T] {
	stages = append(stages, NewStage(stopStage[T]))

	for i := 0; i < len(stages)-1; i++ {
		stages[i].nextFn = stages[i+1].do
	}

	return Pipeline[T]{
		stages: stages,
	}
}

// Do runs the pipeline with the given input and returns the result.
func (p *Pipeline[T]) Do(in T) (T, error) {
	return p.stages[0].do(in)
}

// Do creates pipline and runs the pipeline with the given input and returns the result.
func Do[T any](in T, stages ...StageFn[T]) (T, error) {
	sts := make([]Stage[T], len(stages))
	for i, s := range stages {
		sts[i] = NewStage(s)
	}

	p := New(sts...)

	return p.Do(in)
}

// Run runs the pipeline with default value input.
func (p *Pipeline[T]) Run() (T, error) {
	var v T

	return p.stages[0].do(v)
}

// Run creates pipline and runs the pipeline with default value input.
func Run[T any](stages ...StageFn[T]) (T, error) {
	sts := make([]Stage[T], len(stages))
	for i, s := range stages {
		sts[i] = NewStage(s)
	}

	p := New(sts...)

	return p.Run()
}
