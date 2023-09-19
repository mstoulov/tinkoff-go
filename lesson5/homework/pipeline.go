package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func EnableContextClosing(ctx context.Context, in In) Out {
	out := make(chan any)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				for _ = range in {
				}
				return
			case a, ok := <-in:
				if !ok {
					close(out)
					return
				}
				out <- a
			}
		}
	}()
	return out
}

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	for _, stage := range stages {
		out := stage(EnableContextClosing(ctx, in))
		in = out
	}
	return in
}
