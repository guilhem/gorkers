package gorkers

import "context"

func StopWhenError(ctx context.Context, in interface{}, err error) error {
	return err
}
