package pipeline_test

import (
	"fmt"

	"github.com/burik666/pipeline"
)

func ExamplePipeline() {
	// Define a stage that increments the input by 1
	stageInc := func(in int, next func(int) (int, error)) (int, error) {
		return next(in + 1)
	}

	// Create a new pipeline with two increment stages
	p := pipeline.New(
		pipeline.NewStage(stageInc),
		pipeline.NewStage(stageInc),
	)

	// Run the pipeline with an initial value of 0
	res, err := p.Do(0)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	// Output:
	// 2
}

func ExamplePipeline_producer() {
	// Define a producer that generates numbers from 0 to 4
	producer := func(next func(int) (int, error)) error {
		for i := 0; i < 5; i++ {
			res, err := next(i)
			if err != nil {
				return err
			}

			fmt.Println(res)
		}

		return nil
	}

	// Define a stage that multiplies the input by 2
	stageMul2 := func(in int, next func(int) (int, error)) (int, error) {
		return next(in * 2)
	}

	// Create a new pipeline with a producer and a stage
	p := pipeline.New(
		pipeline.NewProducer(producer),
		pipeline.NewStage(stageMul2),
	)

	// Run the pipeline
	_, err := p.Run()
	if err != nil {
		panic(err)
	}
	// Output:
	// 0
	// 2
	// 4
	// 6
	// 8
}

func ExamplePipeline_Do() {
	// Define a stage that increments the input by 1
	stageInc := func(in int, next func(int) (int, error)) (int, error) {
		return next(in + 1)
	}

	// Run the pipeline with two increment stages
	res, err := pipeline.Do(
		5,
		stageInc,
		stageInc,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	// Output:
	// 7
}

func ExamplePipeline_Middleware() {
	// Define a stage that increments the input by 1
	stageInc := func(in int, next func(int) (int, error)) (int, error) {
		return next(in + 1)
	}

	// Create a new pipeline with two increment stages
	p := pipeline.New(
		pipeline.NewStage(stageInc, pipeline.WithName("stage1")),
		pipeline.NewStage(stageInc, pipeline.WithName("stage2")),
	)

	// Add a middleware that logs the stage name before and after execution
	p.Middleware(func(in int, next func(int) (int, error), opts pipeline.Opts) (int, error) {
		fmt.Printf("pre: %s\n", opts.Name())

		v, err := next(in)

		fmt.Printf("post: %s\n", opts.Name())

		return v, err
	})

	// Run the pipeline with an initial value of 0
	res, err := p.Do(0)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	// Output:
	// pre: stage1
	// pre: stage2
	// post: stage2
	// post: stage1
	// 2
}
