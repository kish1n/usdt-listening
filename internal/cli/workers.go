package cli

import (
	"context"
	"sync"

	"github.com/kish1n/usdt_listening/internal/config"
	"github.com/kish1n/usdt_listening/internal/service/workers"
)

func runServices(ctx context.Context, cfg config.Config, wg *sync.WaitGroup) {
	workers.ListenForTransfers(ctx, cfg)

	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	run(func() { workers.ListenForTransfers(ctx, cfg) })
}
