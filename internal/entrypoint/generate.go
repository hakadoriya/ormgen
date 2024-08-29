package entrypoint

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/hakadoriya/ormgen/internal/apperr"
	"github.com/hakadoriya/ormgen/internal/config"
	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/cliz/clicorez"
	"github.com/hakadoriya/z.go/contextz"
)

func Generate(c *cliz.Command, args []string) error {
	slog.Debug("cmd = "+strings.Join(c.GetExecutedCommandNames(), " "), slog.Any("args", args))
	if len(args) != 1 {
		apperr.Log.ErrorContext(c.Context(), "invalid number of arguments", slog.Any("args", args))
		return clicorez.ErrHelp
	}

	jsonData, _ := json.Marshal(contextz.MustValue[*config.Generate](c.Context()))

	slog.Debug("config = " + string(jsonData))

	return nil
}
