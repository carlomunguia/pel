package main

import (
	"image/color"
	"pel/apptype"
	"pel/pelcanvas"
	"pel/swatch"
	"pel/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	pelApp := app.New()
	pelWindow := pelApp.NewWindow("pel")

	state := apptype.State{
		BrushColor:     color.NRGBA{255, 255, 255, 255},
		SwatchSelected: 0,
	}

	pelCanvasConfig := apptype.PelCanvasConfig{
		DrawingArea:  fyne.NewSize(600, 600),
		CanvasOffset: fyne.NewPos(0, 0),
		PxRows:       10,
		PxCols:       10,
		PxSize:       30,
	}

	pelCanvas := pelcanvas.NewPelCanvas(&state, pelCanvasConfig)

	appInit := ui.AppInit{
		PelCanvas: pelCanvas,
		PelWindow: pelWindow,
		State:     &state,
		Swatches:  make([]*swatch.Swatch, 0, 64),
	}

	ui.Setup(&appInit)

	appInit.PelWindow.ShowAndRun()
}
