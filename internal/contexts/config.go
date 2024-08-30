package contexts

import (
	"context"

	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/z.go/contextz"
)

func GenerateConfig(ctx context.Context) *config.GenerateConfig {
	return contextz.MustValue[*config.GenerateConfig](ctx)
}

func WithGenerateConfig(ctx context.Context, cfg *config.GenerateConfig) context.Context {
	return contextz.WithValue(ctx, *cfg)
}