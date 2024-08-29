package entrypoint

import (
	"github.com/hakadoriya/z.go/buildinfoz"
	"github.com/hakadoriya/z.go/cliz"
)

func Version(c *cliz.Command, _ []string) error {
	return buildinfoz.Fprint(c.Stdout())
}
