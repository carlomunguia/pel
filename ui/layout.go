package ui

func Setup(app *AppInit) {
	swatchesContainer := BuildSwatches(app)

	app.PelWindow.SetContent(swatchesContainer)
}
