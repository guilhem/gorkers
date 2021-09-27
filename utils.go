package gorkers

import "context"

func StopWhenError[I any](ctx context.Context, in I, err error) error {
	return err
}
