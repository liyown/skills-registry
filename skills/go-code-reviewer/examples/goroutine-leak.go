package leak

// Bug: SpawnWorkers returns as soon as it has started the goroutines.
// There is no WaitGroup and no errgroup, so the caller continues
// while the workers may still be running, and there is no signal
// that lets the workers exit. The cancel-on-context branch inside
// the worker is correct on its own, but the workers can still be
// alive when this function returns — the leak is the missing wait
// at the call site, not the context handling inside the worker.

import "context"

func SpawnWorkers(ctx context.Context, jobs <-chan int) {
	for i := 0; i < 5; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-jobs:
					if !ok {
						return
					}
					_ = process(j)
				}
			}
		}()
	}
}

func process(j int) error { return nil }
