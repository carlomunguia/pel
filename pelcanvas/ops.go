// Package pelcanvas provides canvas operations for zooming and panning.
package pelcanvas

import (
	"fyne.io/fyne/v2"
)

// Canvas operation constants
const (
	MinPixelSize     = 2   // Minimum size for a pixel in screen pixels
	MaxPixelSize     = 100 // Maximum size for a pixel in screen pixels
	DefaultPixelSize = 10  // Default pixel size when reset
	ScaleStep        = 1   // Amount to change pixel size per scroll event
)

// scale adjusts the pixel size based on scroll direction
// direction > 0: zoom in (increase pixel size)
// direction < 0: zoom out (decrease pixel size)
// direction = 0: reset to default size
func (pelCanvas *PelCanvas) scale(direction int) {
	oldSize := pelCanvas.PxSize

	switch {
	case direction > 0:
		// Zoom in
		if pelCanvas.PxSize < MaxPixelSize {
			pelCanvas.PxSize += ScaleStep
		}
	case direction < 0:
		// Zoom out
		if pelCanvas.PxSize > MinPixelSize {
			pelCanvas.PxSize -= ScaleStep
		}
	default:
		// Reset to default
		pelCanvas.PxSize = DefaultPixelSize
	}

	// Log if zoom was clamped
	if oldSize == pelCanvas.PxSize && direction != 0 {
		// Zoom was attempted but hit min/max limit
		return
	}
}

// Pan shifts the canvas offset based on mouse movement
// This allows the user to move the canvas around the viewport
func (pelCanvas *PelCanvas) Pan(previousCoord, currentCoord fyne.PointEvent) {
	// Calculate movement delta
	xDiff := currentCoord.Position.X - previousCoord.Position.X
	yDiff := currentCoord.Position.Y - previousCoord.Position.Y

	// Update canvas offset
	pelCanvas.CanvasOffset.X += xDiff
	pelCanvas.CanvasOffset.Y += yDiff

	// Note: Refresh is handled by the caller (MouseMoved) to avoid duplicate refreshes
}

// ResetView resets the canvas to default zoom and position
func (pelCanvas *PelCanvas) ResetView() {
	pelCanvas.PxSize = DefaultPixelSize
	pelCanvas.CanvasOffset = fyne.NewPos(0, 0)
	pelCanvas.Refresh()
}

// ZoomIn increases the pixel size by one step
func (pelCanvas *PelCanvas) ZoomIn() {
	pelCanvas.scale(1)
	pelCanvas.Refresh()
}

// ZoomOut decreases the pixel size by one step
func (pelCanvas *PelCanvas) ZoomOut() {
	pelCanvas.scale(-1)
	pelCanvas.Refresh()
}

// ZoomToFit adjusts the canvas to fit the window size
func (pelCanvas *PelCanvas) ZoomToFit(windowSize fyne.Size) {
	// Calculate optimal pixel size to fit canvas in window
	maxPxSizeWidth := int(windowSize.Width / float32(pelCanvas.PxCols))
	maxPxSizeHeight := int(windowSize.Height / float32(pelCanvas.PxRows))

	// Use the smaller of the two to ensure entire canvas fits
	optimalSize := maxPxSizeWidth
	if maxPxSizeHeight < maxPxSizeWidth {
		optimalSize = maxPxSizeHeight
	}

	// Clamp to valid range
	if optimalSize < MinPixelSize {
		optimalSize = MinPixelSize
	}
	if optimalSize > MaxPixelSize {
		optimalSize = MaxPixelSize
	}

	pelCanvas.PxSize = optimalSize
	pelCanvas.CanvasOffset = fyne.NewPos(0, 0)
	pelCanvas.Refresh()
}

// GetZoomLevel returns the current zoom level as a percentage
// 100% = DefaultPixelSize
func (pelCanvas *PelCanvas) GetZoomLevel() int {
	return (pelCanvas.PxSize * 100) / DefaultPixelSize
}
