package utility

import (
	"context"
	"math"
	"time"
)

type ExponentialBackoff struct {
	attempts int
	factor   float64
	min      time.Duration
	max      time.Duration
}

func NewExponentialBackoff(
	factor float64,
	min time.Duration,
	max time.Duration,
) *ExponentialBackoff {
	return &ExponentialBackoff{
		factor: factor,
		min:    min,
		max:    max,
	}
}

func (b *ExponentialBackoff) Next() time.Duration {
	b.attempts++
	duration := float64(b.min) * math.Pow(b.factor, float64(b.attempts-1))
	if duration > float64(b.max) {
		return b.max
	}
	return time.Duration(duration)
}

func RetryWithBackoff(ctx context.Context, maxAttempts int, backoff *ExponentialBackoff, operation func() error) error {
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := operation(); err != nil {
			lastErr = err
			delay := backoff.Next()

			SugarLogger.Debugw("retry attempt failed",
				"attempt", i+1,
				"max_attempts", maxAttempts,
				"delay", delay,
				"error", err,
			)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				continue
			}
		}
		return nil
	}

	return lastErr
}
