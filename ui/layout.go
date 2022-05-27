package ui

import "fyne.io/fyne/v2/container"

func Setup(app *AppInit) {
	SetupMenus(app)
	swatchesContainer := BuildSwatches(app)
	colorPicker := SetupColorPicker(app)

	appLayout := container.NewBorder(nil, swatchesContainer, nil, colorPicker, app.PelCanvas)

	app.PelWindow.SetContent(appLayout)
}
