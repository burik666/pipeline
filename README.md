# Pipeline

A lightweight and flexible generic pipeline library for Go.
It allows you to define and compose stages and producers that transform or generate data in a pipeline style.

## Installation

```bash
go get github.com/burik666/pipeline
```

## Overview

The package provides the following core concepts:

- **StageFn** — a function that transforms input and optionally calls the next stage:
  ```go
  func(in T, next func(T) (T, error)) (T, error)
  ```
- **NewStage** — creates a stage with access to `next`.
- **NewSimpleStage** — a stage without `next`, useful for final transformations.
- **NewProducer** — a stage that emits values instead of consuming input.
- **Pipeline** — a sequence of stages, created with `New` and executed via `Do` (with input) or `Run` (without explicit input).
- **Middleware** — intercepts execution of stages (for logging, tracing, metrics).
- **WithName** — attaches names to stages for easier debugging.
- **Reverse result propagation** — stages can return values back up the chain instead of (or in addition to) passing them forward.
  This allows implementing finalizers, cleanup logic, or collecting results after downstream processing.

---

## Usage

### Simple Pipeline

This example shows the most basic usage of a pipeline:
we define a stage that increments the input and chain it twice.
The input value `5` goes through two stages (`+1`, `+1`) and produces `7`.

```go
package main

import (
    "fmt"
    "github.com/burik666/pipeline"
)

func main() {
    inc := func(in int, next func(int) (int, error)) (int, error) {
        return next(in + 1)
    }

    result, err := pipeline.Do(
        5,
        inc,
        inc,
    )
    if err != nil {
        panic(err)
    }

    fmt.Println(result)
    // Output: 7
}
```

---

### Producer Example

Here we demonstrate a **producer stage**:
instead of receiving input, it generates a sequence of numbers (0 to 4) and sends them down the pipeline.
Each value is then multiplied by 2 in the next stage.
The final pipeline prints `0, 2, 4, 6, 8`.

```go
package main

import (
    "fmt"
    "github.com/burik666/pipeline"
)

func main() {
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

    mul2 := func(in int, next func(int) (int, error)) (int, error) {
        return next(in * 2)
    }

    p := pipeline.New(
        pipeline.NewProducer(producer),
        pipeline.NewStage(mul2),
    )

    _, err := p.Run()
    if err != nil {
        panic(err)
    }
}

// Output:
// 0
// 2
// 4
// 6
// 8
```

---

### Middleware Example

This example shows how **middleware** can wrap stages to perform additional logic.
We add names to stages (`stage1`, `stage2`) and use middleware to log messages before and after each stage.
This demonstrates how you can implement logging, tracing, or metrics around every pipeline step.

```go
package main

import (
    "fmt"
    "github.com/burik666/pipeline"
)

func main() {
    inc := func(in int, next func(int) (int, error)) (int, error) {
        return next(in + 1)
    }

    p := pipeline.New(
        pipeline.NewStage(inc, pipeline.WithName("stage1")),
        pipeline.NewStage(inc, pipeline.WithName("stage2")),
    )

    p.Middleware(func(in int, next func(int) (int, error), opts pipeline.Opts) (int, error) {
        fmt.Printf("pre: %s\n", opts.Name())
        v, err := next(in)
        fmt.Printf("post: %s\n", opts.Name())
        return v, err
    })

    res, err := p.Do(0)
    if err != nil {
        panic(err)
    }

    fmt.Println(res)
}

// Output:
// pre: stage1
// pre: stage2
// post: stage2
// post: stage1
// 2
```

---

### Simple composition

This pipeline demonstrates how to **combine multiple transformations**:
first, the number is doubled, then `3` is added.
Input `4` goes through `*2 → +3` and produces `11`.

```go
func ExampleDo_simplePipeline() {
    double := func(in int, next func(int) (int, error)) (int, error) {
        return next(in * 2)
    }
    addThree := func(in int, next func(int) (int, error)) (int, error) {
        return next(in + 3)
    }

    result, err := pipeline.Do(4, double, addThree)
    if err != nil {
        panic(err)
    }

    fmt.Println(result)
    // Output: 11 (4 * 2 + 3)
}
```

---

### Producer with `Run`

This test shows how to use a **producer together with `Run`**.
The producer emits a single value (`1`), which is then transformed by the next stage (`+5`).
The result of running the pipeline is `6`.

```go
func ExampleRun_withDefault() {
    produceOne := func(next func(int) (int, error)) error {
        next(1)
        return nil
    }
    addFive := func(in int, next func(int) (int, error)) (int, error) {
        return next(in + 5)
    }

    p := pipeline.New(
        pipeline.NewProducer(produceOne),
        pipeline.NewStage(addFive),
    )

    result, err := p.Run()
    if err != nil {
        panic(err)
    }

    fmt.Println(result)
    // Output: 6
}
```

---

## License

MIT License — see the [LICENSE](LICENSE) file for details.
