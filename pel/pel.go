package main

import (
	"image/color"
	"pel/apptype"
	"pel/swatch"
	"pel/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	pelApp := app.New()
	pelWindow := pelApp.NewWindow("pel")

	state := apptype.State{
		BrushColor:     color.NRGBA{255, 255, 255, 255},
		SwatchSelected: 0,
	}

	appInit := ui.AppInit{
		PelWindow: pelWindow,
		State:     &state,
		Swatches:  make([]*swatch.Swatch, 0, 64),
	}

	ui.Setup(&appInit)

	appInit.PelWindow.ShowAndRun()
}
