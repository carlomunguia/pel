// Package pelcanvas provides rendering functionality for the pixel canvas widget.
package pelcanvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Border position indices
const (
	BorderLeft   = 0
	BorderTop    = 1
	BorderRight  = 2
	BorderBottom = 3
	BorderCount  = 4
)

// PelCanvasRenderer handles the rendering of the pixel canvas widget
type PelCanvasRenderer struct {
	pelCanvas    *PelCanvas          // Reference to the canvas widget
	canvasImage  *canvas.Image       // The pixel art image being displayed
	canvasBorder []canvas.Line       // Border lines around the canvas
	canvasCursor []fyne.CanvasObject // Current cursor objects
}

// MinSize returns the minimum size required to display the canvas
func (renderer *PelCanvasRenderer) MinSize() fyne.Size {
	if renderer.pelCanvas == nil {
		return fyne.NewSize(0, 0)
	}
	return renderer.pelCanvas.DrawingArea
}

// Objects returns all canvas objects that need to be rendered
func (renderer *PelCanvasRenderer) Objects() []fyne.CanvasObject {
	// Pre-allocate with exact capacity: borders + image + cursor objects
	capacity := len(renderer.canvasBorder) + 1 + len(renderer.canvasCursor)
	objects := make([]fyne.CanvasObject, 0, capacity)

	// Add border lines
	for i := range renderer.canvasBorder {
		objects = append(objects, &renderer.canvasBorder[i])
	}

	// Add canvas image
	if renderer.canvasImage != nil {
		objects = append(objects, renderer.canvasImage)
	}

	// Add cursor objects
	objects = append(objects, renderer.canvasCursor...)

	return objects
}

// Destroy cleans up any resources used by the renderer
func (renderer *PelCanvasRenderer) Destroy() {
	// Clean up resources
	renderer.canvasImage = nil
	renderer.canvasBorder = nil
	renderer.canvasCursor = nil
	renderer.pelCanvas = nil
}

// Layout positions and sizes all canvas objects
func (renderer *PelCanvasRenderer) Layout(size fyne.Size) {
	if renderer.pelCanvas == nil {
		return
	}

	renderer.layoutCanvas(size)
	renderer.layoutBorder(size)
}

// layoutCanvas positions and sizes the main canvas image
func (renderer *PelCanvasRenderer) layoutCanvas(size fyne.Size) {
	if renderer.canvasImage == nil {
		return
	}

	imgPxWidth := renderer.pelCanvas.PxCols
	imgPxHeight := renderer.pelCanvas.PxRows
	pxSize := renderer.pelCanvas.PxSize

	// Position the image at the canvas offset
	renderer.canvasImage.Move(fyne.NewPos(
		renderer.pelCanvas.CanvasOffset.X,
		renderer.pelCanvas.CanvasOffset.Y,
	))

	// Size the image based on pixel grid dimensions
	renderer.canvasImage.Resize(fyne.NewSize(
		float32(imgPxWidth*pxSize),
		float32(imgPxHeight*pxSize),
	))
}

// layoutBorder positions the border lines around the canvas
func (renderer *PelCanvasRenderer) layoutBorder(size fyne.Size) {
	if len(renderer.canvasBorder) < BorderCount {
		return
	}

	if renderer.canvasImage == nil {
		return
	}

	offset := renderer.pelCanvas.CanvasOffset
	imgHeight := renderer.canvasImage.Size().Height
	imgWidth := renderer.canvasImage.Size().Width

	// Left border
	left := &renderer.canvasBorder[BorderLeft]
	left.Position1 = fyne.NewPos(offset.X, offset.Y)
	left.Position2 = fyne.NewPos(offset.X, offset.Y+imgHeight)

	// Top border
	top := &renderer.canvasBorder[BorderTop]
	top.Position1 = fyne.NewPos(offset.X, offset.Y)
	top.Position2 = fyne.NewPos(offset.X+imgWidth, offset.Y)

	// Right border
	right := &renderer.canvasBorder[BorderRight]
	right.Position1 = fyne.NewPos(offset.X+imgWidth, offset.Y)
	right.Position2 = fyne.NewPos(offset.X+imgWidth, offset.Y+imgHeight)

	// Bottom border
	bottom := &renderer.canvasBorder[BorderBottom]
	bottom.Position1 = fyne.NewPos(offset.X, offset.Y+imgHeight)
	bottom.Position2 = fyne.NewPos(offset.X+imgWidth, offset.Y+imgHeight)
}

// Refresh updates the renderer with the latest canvas state
func (renderer *PelCanvasRenderer) Refresh() {
	if renderer.pelCanvas == nil {
		return
	}

	// Reload image if needed (e.g., after loading a new file)
	if renderer.pelCanvas.reloadImage {
		// Clean up old image
		oldImage := renderer.canvasImage

		// Create new image
		renderer.canvasImage = canvas.NewImageFromImage(renderer.pelCanvas.PixelData)
		renderer.canvasImage.ScaleMode = canvas.ImageScalePixels
		renderer.canvasImage.FillMode = canvas.ImageFillContain

		// Clear reload flag
		renderer.pelCanvas.reloadImage = false

		// Allow old image to be garbage collected
		_ = oldImage
	}

	// Update layout and refresh image
	renderer.Layout(renderer.pelCanvas.Size())

	if renderer.canvasImage != nil {
		canvas.Refresh(renderer.canvasImage)
	}
}

// SetCursor updates the cursor objects to be displayed
func (renderer *PelCanvasRenderer) SetCursor(objects []fyne.CanvasObject) {
	if objects == nil {
		renderer.canvasCursor = make([]fyne.CanvasObject, 0)
	} else {
		renderer.canvasCursor = objects
	}
}
