package entrypoint

import (
	"github.com/hakadoriya/ormgen/internal/consts"
	"github.com/hakadoriya/ormgen/internal/contexts"
	source_go "github.com/hakadoriya/ormgen/internal/lang/go/source"
	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/errorz"
)

func Generate(c *cliz.Command, args []string) error {
	cfg := contexts.GenerateConfig(c.Context())

	switch cfg.Language {
	case consts.LanguageGo:
		if err := source_go.Parse(c.Context(), args); err != nil {
			return errorz.Errorf("parsego.Parse: %w", err)
		}
	}

	return nil
}
