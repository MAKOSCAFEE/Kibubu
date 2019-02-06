package grifts

import (
	"github.com/barnie/kibubu/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
