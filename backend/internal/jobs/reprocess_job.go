package jobs

import "context"

func EnqueueReprocess(_ context.Context, _ string) error {
	// Hook for background worker queue implementation.
	return nil
}
