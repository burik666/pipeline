package pipeline

// Middleware adds middleware to the pipeline.
// The middleware is applied to each stage in the pipeline.
// The middleware function is called with the input, the next function, and the stage options.
func (p *Pipeline[T]) Middleware(mw func(T, func(T) (T, error), Opts) (T, error)) {
	// skip the last stop stage
	for i := 0; i < len(p.stages)-1; i++ {
		st := p.stages[i]
		p.stages[i] = NewStage(func(in T, next func(T) (T, error)) (T, error) {
			return mw(in, next, st.opts)
		},
			withOpts(st.opts),
		)
		p.stages[i].nextFn = st.do
	}
}
