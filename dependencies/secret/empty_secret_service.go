package secret

import (
	"context"
	"github.com/InfluxCommunity/flux/codes"
	"github.com/InfluxCommunity/flux/internal/errors"
)

func (ess EmptySecretService) LoadSecret(ctx context.Context, k string) (string, error) {
	return "", errors.Newf(codes.NotFound, "secret key %q not found", k)
}

// Secret service that always reports no secrets exist
type EmptySecretService struct {
}
