package entrypoint

import (
	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	gen_go "github.com/hakadoriya/ormgen/internal/lang/go/gen"
	source_go "github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/ormgen/pkg/apperr"
	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/errorz"
)

func Generate(c *cliz.Command, args []string) error {
	cfg := contexts.GenerateConfig(c.Context())

	switch cfg.Language {
	case consts.LanguageGo:
		packageSources, err := source_go.Parse(c.Context(), args)
		if err != nil {
			return errorz.Errorf("parsego.Parse: %w", err)
		}

		if err := gen_go.Output(c.Context(), packageSources); err != nil {
			return errorz.Errorf("gen.Output: %w", err)
		}

	default:
		return errorz.Errorf("lang=%s: %w", cfg.Language, apperr.ErrLanguageNotSupported)
	}

	return nil
}
