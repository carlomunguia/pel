package ui

import (
	"pel/apptype"
	"pel/swatch"

	"fyne.io/fyne/v2"
)

type AppInit struct {
	PelWindow fyne.Window
	State     *apptype.State
	Swatches  []*swatch.Swatch
}
