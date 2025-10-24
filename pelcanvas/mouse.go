// Package pelcanvas provides mouse event handling for the pixel canvas.
package pelcanvas

import (
	"github.com/carlomunguia/pel/pelcanvas/brush"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// Mouse interaction constants
const (
	ScrollSensitivity = 1 // Multiplier for scroll events
)

// Scrolled handles mouse scroll events for zooming the canvas
func (pelCanvas *PelCanvas) Scrolled(ev *fyne.ScrollEvent) {
	if ev == nil {
		return
	}

	// Scale canvas based on scroll direction
	scrollDelta := int(ev.Scrolled.DY * ScrollSensitivity)
	pelCanvas.scale(scrollDelta)
	pelCanvas.Refresh()
}

// MouseMoved handles mouse movement events for drawing and cursor updates
func (pelCanvas *PelCanvas) MouseMoved(ev *desktop.MouseEvent) {
	if ev == nil {
		return
	}

	// Update cursor and handle drawing
	needsRefresh := pelCanvas.updateCursorAndDraw(ev)

	// Handle canvas panning
	pelCanvas.TryPan(pelCanvas.mouseState.previousCoord, ev)

	// Store current position for next event
	pelCanvas.mouseState.previousCoord = &ev.PointEvent

	// Single refresh at the end if needed
	if needsRefresh {
		pelCanvas.Refresh()
	}
}

// MouseDown handles mouse button press events
func (pelCanvas *PelCanvas) MouseDown(ev *desktop.MouseEvent) {
	if ev == nil {
		return
	}

	// Attempt to draw/interact with brush
	if brush.TryBrush(pelCanvas.appState, pelCanvas, ev) {
		pelCanvas.Refresh()
	}
}

// MouseIn handles mouse entering the canvas area
func (pelCanvas *PelCanvas) MouseIn(ev *desktop.MouseEvent) {
	// TODO: Implement mouse enter behavior (e.g., show cursor, enable drawing)
}

// MouseUp handles mouse button release events
func (pelCanvas *PelCanvas) MouseUp(ev *desktop.MouseEvent) {
	// TODO: Implement mouse up behavior (e.g., finish line/shape drawing)
}

// MouseOut handles mouse leaving the canvas area
func (pelCanvas *PelCanvas) MouseOut() {
	// Clear cursor when mouse leaves canvas
	if pelCanvas.renderer != nil {
		pelCanvas.renderer.SetCursor(make([]fyne.CanvasObject, 0))
		pelCanvas.Refresh()
	}
}

// updateCursorAndDraw updates the cursor position and handles drawing operations
// Returns true if a refresh is needed
func (pelCanvas *PelCanvas) updateCursorAndDraw(ev *desktop.MouseEvent) bool {
	x, y := pelCanvas.MouseToCanvasXY(ev)

	if x != nil && y != nil {
		// Mouse is over a valid canvas pixel

		// Try to draw if mouse button is pressed (handled by TryBrush)
		drawn := brush.TryBrush(pelCanvas.appState, pelCanvas, ev)

		// Update cursor to show current brush tool
		cursor := brush.Cursor(
			pelCanvas.PelCanvasConfig,
			pelCanvas.appState.BrushType,
			ev,
			*x,
			*y,
		)

		if pelCanvas.renderer != nil {
			pelCanvas.renderer.SetCursor(cursor)
		}

		return drawn // Only refresh if something was drawn
	} else {
		// Mouse is outside valid canvas area - hide cursor
		if pelCanvas.renderer != nil {
			pelCanvas.renderer.SetCursor(make([]fyne.CanvasObject, 0))
		}
		return true // Refresh to clear cursor
	}
}
