package pkg

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotReadyYet        = errors.New("not ready yet")
	ErrWaitTimeExceeded   = errors.New("wait time exceeded")
	ErrOtherErrorOccurred = errors.New("other error occurred")
)

type Waiter struct {
	ticker   *time.Ticker
	numTicks uint8
}

func NewWaiter(tick time.Duration, numTicks uint8) *Waiter {
	return &Waiter{ticker: time.NewTicker(tick), numTicks: numTicks}
}

func (w *Waiter) Wait(ctx context.Context, f func(ctx context.Context) error) error {
	for ticksPassed := uint8(0); ticksPassed < w.numTicks; ticksPassed += 1 {
		select {
		case <-w.ticker.C:
			err := f(ctx)
			if err == nil {
				return nil
			}

			if !errors.Is(err, ErrNotReadyYet) {
				return fmt.Errorf("%s: %s", ErrOtherErrorOccurred.Error(), err.Error())
			}
		case <-ctx.Done():
			return nil
		}
	}

	return ErrWaitTimeExceeded
}
