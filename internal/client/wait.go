package client

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func Ticker(tick time.Duration) wait.WaitWithContextFunc {
	return func(ctx context.Context) <-chan struct{} {
		ticker := make(chan struct{})

		go func() {
			defer close(ticker)
			for {
				select {
				case <-time.After(tick):
					ticker <- struct{}{}
				case <-ctx.Done():
					return
				}
			}
		}()

		return ticker
	}
}
