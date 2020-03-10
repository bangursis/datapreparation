package profiles

import "context"

type Repository interface {
	Save(ctx context.Context, iccID string, encrypted [][]byte) error
}
