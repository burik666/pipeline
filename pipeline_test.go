package pipeline_test

import (
	"testing"

	"github.com/burik666/pipeline"
)

func TestPipeline(t *testing.T) {
	stageInc := func(in int, next func(int) (int, error)) (int, error) {
		return next(in + 1)
	}

	p := pipeline.New(
		pipeline.NewStage(stageInc, pipeline.WithName("stage1")),
		pipeline.NewStage(stageInc, pipeline.WithName("stage2")),
		pipeline.NewStage(stageInc, pipeline.WithName("stage3")),
	)

	p.Middleware(func(in int, next func(int) (int, error), _ pipeline.Opts) (int, error) {
		return next(in)
	})

	// Do
	res, err := p.Do(0)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want := 3
	if res != want {
		t.Errorf("unexpected value: got: %d: want: %d", res, want)
	}
}

func TestPipelineRun(t *testing.T) {
	_, err := pipeline.Run(func(in int, next func(int) (int, error)) (int, error) {
		return next(in)
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func BenchmarkPipeline(b *testing.B) {
	fn := func(in int, next func(int) (int, error)) (int, error) {
		return next(in)
	}

	p := pipeline.New(
		pipeline.NewStage(fn, pipeline.WithName("stage1")),
		pipeline.NewStage(fn, pipeline.WithName("stage2")),
		pipeline.NewStage(fn, pipeline.WithName("stage3")),
	)

	p.Middleware(func(in int, next func(int) (int, error), _ pipeline.Opts) (int, error) {
		return next(in)
	})

	for i := 0; i < b.N; i++ {
		_, err := p.Do(i)
		if err != nil {
			b.Errorf("unexpected error: %s", err)
		}
	}
}
