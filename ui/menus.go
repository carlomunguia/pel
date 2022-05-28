package ui

import (
	"errors"
	"image"
	"image/png"
	"os"
	"pel/util"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func saveFileDialog(app *AppInit) {
	dialog.ShowFileSave(func(uri fyne.URIWriteCloser, e error) {
		if uri == nil {
			return
		} else {
			err := png.Encode(uri, app.PelCanvas.PixelData)
			if err != nil {
				dialog.ShowError(err, app.PelWindow)
				return
			}
			app.State.SetFilePath(uri.URI().Path())
		}
	}, app.PelWindow)
}

func BuildMenus(app *AppInit) *fyne.Menu {
	return fyne.NewMenu(
		"File",
		BuildNewMenu(app),
		BuildOpenMenu(app),
		BuildSaveMenu(app),
		BuildSaveAsMenu(app),
	)
}

func BuildNewMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("New", func() {
		sizeValidator := func(s string) error {
			width, err := strconv.Atoi(s)
			if err != nil {
				return errors.New("must be a positive integer")
			}
			if width <= 0 {
				return errors.New("must be greater than zero")
			}
			return nil
		}
		widthEntry := widget.NewEntry()
		widthEntry.Validator = sizeValidator

		heightEntry := widget.NewEntry()
		heightEntry.Validator = sizeValidator

		widthFormEntry := widget.NewFormItem("Width", widthEntry)
		heightFormEntry := widget.NewFormItem("Height", heightEntry)

		formItems := []*widget.FormItem{widthFormEntry, heightFormEntry}

		dialog.ShowForm("New Image", "Create", "Cancel", formItems, func(ok bool) {
			if ok {
				pixelWidth := 0
				pixelHeight := 0
				if widthEntry.Validate() != nil {
					dialog.ShowError(errors.New("invalid width"), app.PelWindow)
				} else {
					pixelWidth, _ = strconv.Atoi(widthEntry.Text)
				}
				if heightEntry.Validate() != nil {
					dialog.ShowError(errors.New("invalid height"), app.PelWindow)
				} else {
					pixelHeight, _ = strconv.Atoi(heightEntry.Text)
				}
				app.PelCanvas.NewDrawing(pixelWidth, pixelHeight)
			}
		}, app.PelWindow)
	})
}

func BuildOpenMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Open", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri == nil {
				return
			} else {
				file, _, err := image.Decode(uri)
				if err != nil {
					dialog.ShowError(err, app.PelWindow)
					return
				}
				app.PelCanvas.LoadImage(file)
				app.State.SetFilePath(uri.URI().Path())
				imgColors := util.GetImageColors(file)
				i := 0
				for c := range imgColors {
					if i == len(app.Swatches) {
						break
					}
					app.Swatches[i].SetColor(c)
					i++
				}
			}
		}, app.PelWindow)
	})
}

func BuildSaveMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Save", func() {
		if app.State.FilePath == "" {
			saveFileDialog(app)
		} else {
			tryClose := func(fh *os.File) {
				err := fh.Close()
				if err != nil {
					dialog.ShowError(err, app.PelWindow)
				}
			}
			fh, err := os.Create(app.State.FilePath)
			defer tryClose(fh)

			if err != nil {
				dialog.ShowError(err, app.PelWindow)
				return
			}
			err = png.Encode(fh, app.PelCanvas.PixelData)
			if err != nil {
				dialog.ShowError(err, app.PelWindow)
				return
			}
		}
	})
}

func BuildSaveAsMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Save As...", func() {
		saveFileDialog(app)
	})
}

func SetupMenus(app *AppInit) {
	menus := BuildMenus(app)
	mainMenu := fyne.NewMainMenu(menus)
	app.PelWindow.SetMainMenu(mainMenu)
}
