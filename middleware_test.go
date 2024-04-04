package pipeline_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/burik666/pipeline"
)

func TestMiddleware(t *testing.T) {
	type Data struct {
		Steps []string
	}

	stageInc := func(in Data, next func(Data) (Data, error)) (Data, error) {
		in.Steps = append(in.Steps, "inc")

		return next(in)
	}

	p := pipeline.New(
		pipeline.NewStage(stageInc, pipeline.WithName("stage1")),
		pipeline.NewStage(stageInc, pipeline.WithName("stage2")),
	)

	p.Middleware(func(in Data, next func(Data) (Data, error), opts pipeline.Opts) (Data, error) {
		in.Steps = append(in.Steps, fmt.Sprintf("mw1 pre: %s", opts.Name()))
		v, err := next(in)
		v.Steps = append(v.Steps, fmt.Sprintf("mw1 post: %s", opts.Name()))

		return v, err
	})

	p.Middleware(func(in Data, next func(Data) (Data, error), opts pipeline.Opts) (Data, error) {
		in.Steps = append(in.Steps, fmt.Sprintf("mw2 pre: %s", opts.Name()))
		v, err := next(in)
		v.Steps = append(v.Steps, fmt.Sprintf("mw2 post: %s", opts.Name()))

		return v, err
	})

	res, _ := p.Do(Data{})

	if !slices.Equal(res.Steps, []string{
		"mw2 pre: stage1",
		"mw1 pre: stage1",
		"inc",
		"mw2 pre: stage2",
		"mw1 pre: stage2",
		"inc",
		"mw1 post: stage2",
		"mw2 post: stage2",
		"mw1 post: stage1",
		"mw2 post: stage1",
	}) {
		t.Fatalf("unexpected steps: %#v", res.Steps)
	}
}
