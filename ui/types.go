package ui

import (
	"pel/apptype"
	"pel/pelcanvas"
	"pel/swatch"

	"fyne.io/fyne/v2"
)

type AppInit struct {
	PelCanvas *pelcanvas.PelCanvas
	PelWindow fyne.Window
	State     *apptype.State
	Swatches  []*swatch.Swatch
}
