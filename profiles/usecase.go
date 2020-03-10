package profiles

import "context"

type UseCase interface {
	Import(ctx context.Context, filepath string) ([]byte, error)
}
