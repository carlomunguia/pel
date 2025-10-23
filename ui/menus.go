// Package ui provides menu construction and handling for the Pel pixel art editor.
package ui

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"pel/util"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Menu constants
const (
	DefaultImageWidth  = 32
	DefaultImageHeight = 32
	MinImageSize       = 1
	MaxImageSize       = 1024
)

// Error messages
var (
	ErrInvalidWidth  = errors.New("width must be a positive integer")
	ErrInvalidHeight = errors.New("height must be a positive integer")
	ErrNoFile        = errors.New("no file selected")
	ErrSaveFailed    = errors.New("failed to save file")
)

// SetupMenus initializes and attaches the main menu to the application window
func SetupMenus(app *AppInit) {
	if app == nil || app.PelWindow == nil {
		log.Println("Warning: Cannot setup menus - app or window is nil")
		return
	}

	menus := BuildMenus(app)
	mainMenu := fyne.NewMainMenu(menus)
	app.PelWindow.SetMainMenu(mainMenu)
	log.Println("Menus initialized successfully")
}

// BuildMenus constructs the main application menu structure
func BuildMenus(app *AppInit) *fyne.Menu {
	return fyne.NewMenu(
		"File",
		BuildNewMenu(app),
		BuildOpenMenu(app),
		fyne.NewMenuItemSeparator(),
		BuildSaveMenu(app),
		BuildSaveAsMenu(app),
		fyne.NewMenuItemSeparator(),
		BuildExportMenu(app),
		fyne.NewMenuItemSeparator(),
		BuildQuitMenu(app),
	)
}

// BuildNewMenu creates the "New" menu item for creating a new image
func BuildNewMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("New", func() {
		showNewImageDialog(app)
	})
}

// BuildOpenMenu creates the "Open" menu item for loading an image
func BuildOpenMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Open...", func() {
		showOpenFileDialog(app)
	})
}

// BuildSaveMenu creates the "Save" menu item
func BuildSaveMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Save", func() {
		saveImage(app, false)
	})
}

// BuildSaveAsMenu creates the "Save As" menu item
func BuildSaveAsMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Save As...", func() {
		saveImage(app, true)
	})
}

// BuildExportMenu creates the "Export" menu item for exporting to different formats
func BuildExportMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Export...", func() {
		// TODO: Implement export with format options (JPEG, GIF, etc.)
		dialog.ShowInformation("Export", "Export feature coming soon!", app.PelWindow)
	})
}

// BuildQuitMenu creates the "Quit" menu item
func BuildQuitMenu(app *AppInit) *fyne.MenuItem {
	return fyne.NewMenuItem("Quit", func() {
		app.PelWindow.Close()
	})
}

// showNewImageDialog displays a dialog for creating a new image
func showNewImageDialog(app *AppInit) {
	if app == nil {
		return
	}

	// Create size validator
	sizeValidator := func(s string) error {
		if s == "" {
			return errors.New("size cannot be empty")
		}
		size, err := strconv.Atoi(s)
		if err != nil {
			return ErrInvalidWidth
		}
		if size < MinImageSize {
			return fmt.Errorf("size must be at least %d", MinImageSize)
		}
		if size > MaxImageSize {
			return fmt.Errorf("size must be at most %d", MaxImageSize)
		}
		return nil
	}

	// Create form entries with default values
	widthEntry := widget.NewEntry()
	widthEntry.SetText(strconv.Itoa(DefaultImageWidth))
	widthEntry.Validator = sizeValidator

	heightEntry := widget.NewEntry()
	heightEntry.SetText(strconv.Itoa(DefaultImageHeight))
	heightEntry.Validator = sizeValidator

	widthFormEntry := widget.NewFormItem("Width", widthEntry)
	heightFormEntry := widget.NewFormItem("Height", heightEntry)

	formItems := []*widget.FormItem{widthFormEntry, heightFormEntry}

	dialog.ShowForm("New Image", "Create", "Cancel", formItems, func(ok bool) {
		if !ok {
			return
		}

		// Validate and parse dimensions
		if err := widthEntry.Validate(); err != nil {
			dialog.ShowError(fmt.Errorf("invalid width: %w", err), app.PelWindow)
			return
		}
		if err := heightEntry.Validate(); err != nil {
			dialog.ShowError(fmt.Errorf("invalid height: %w", err), app.PelWindow)
			return
		}

		pixelWidth, _ := strconv.Atoi(widthEntry.Text)
		pixelHeight, _ := strconv.Atoi(heightEntry.Text)

		// Create new drawing
		if err := app.PelCanvas.NewDrawing(pixelWidth, pixelHeight); err != nil {
			dialog.ShowError(fmt.Errorf("failed to create new drawing: %w", err), app.PelWindow)
			return
		}

		log.Printf("Created new image: %dx%d", pixelWidth, pixelHeight)
	}, app.PelWindow)
}

// showOpenFileDialog displays a file open dialog
func showOpenFileDialog(app *AppInit) {
	if app == nil {
		return
	}

	dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to open file: %w", err), app.PelWindow)
			return
		}
		if uri == nil {
			return
		}
		defer uri.Close()

		// Decode image
		img, format, err := image.Decode(uri)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to decode image: %w", err), app.PelWindow)
			return
		}

		log.Printf("Opened image: %s (format: %s)", uri.URI().Path(), format)

		// Load image into canvas
		if err := app.PelCanvas.LoadImage(img); err != nil {
			dialog.ShowError(fmt.Errorf("failed to load image: %w", err), app.PelWindow)
			return
		}

		// Update file path
		app.State.SetFilePath(uri.URI().Path())

		// Extract colors and update swatches
		updateSwatchesFromImage(app, img)

		dialog.ShowInformation("Success",
			fmt.Sprintf("Loaded: %s\nSize: %dx%d",
				filepath.Base(uri.URI().Path()),
				img.Bounds().Dx(),
				img.Bounds().Dy()),
			app.PelWindow)
	}, app.PelWindow)
}

// saveImage saves the current image to disk
func saveImage(app *AppInit, forceDialog bool) {
	if app == nil {
		return
	}

	// Show save dialog if no file path or force dialog
	if app.State.FilePath == "" || forceDialog {
		showSaveFileDialog(app)
		return
	}

	// Save to existing file path
	if err := saveImageToFile(app, app.State.FilePath); err != nil {
		dialog.ShowError(fmt.Errorf("save failed: %w", err), app.PelWindow)
		return
	}

	log.Printf("Saved image to: %s", app.State.FilePath)
	dialog.ShowInformation("Success",
		fmt.Sprintf("Saved: %s", filepath.Base(app.State.FilePath)),
		app.PelWindow)
}

// showSaveFileDialog displays a file save dialog
func showSaveFileDialog(app *AppInit) {
	if app == nil {
		return
	}

	dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to save file: %w", err), app.PelWindow)
			return
		}
		if uri == nil {
			return
		}
		defer uri.Close()

		// Encode and write image
		if err := png.Encode(uri, app.PelCanvas.PixelData); err != nil {
			dialog.ShowError(fmt.Errorf("failed to encode image: %w", err), app.PelWindow)
			return
		}

		// Update file path
		filePath := uri.URI().Path()
		app.State.SetFilePath(filePath)

		log.Printf("Saved image to: %s", filePath)
		dialog.ShowInformation("Success",
			fmt.Sprintf("Saved: %s", filepath.Base(filePath)),
			app.PelWindow)
	}, app.PelWindow)
}

// saveImageToFile saves the image to the specified file path
func saveImageToFile(app *AppInit, filePath string) error {
	if app == nil || app.PelCanvas == nil {
		return errors.New("invalid app state")
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Warning: failed to close file: %v", closeErr)
		}
	}()

	// Encode image
	if err := png.Encode(file, app.PelCanvas.PixelData); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

// updateSwatchesFromImage extracts colors from the image and updates swatches
func updateSwatchesFromImage(app *AppInit, img image.Image) {
	if app == nil || img == nil || len(app.Swatches) == 0 {
		return
	}

	// Use the updated GetImageColors function with error handling
	imgColors, err := util.GetImageColors(img)
	if err != nil {
		log.Printf("Warning: Failed to extract colors: %v", err)
		return
	}

	i := 0
	for c := range imgColors {
		if i >= len(app.Swatches) {
			break
		}
		app.Swatches[i].SetColor(c)
		i++
	}

	log.Printf("Updated %d swatches from image colors", i)
}
