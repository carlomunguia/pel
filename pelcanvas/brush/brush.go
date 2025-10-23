// Package brush provides brush tools and cursor rendering for the Pel pixel art editor.
package brush

import (
	"image/color"
	"pel/apptype"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

// Cursor configuration constants
const (
	DefaultCursorStrokeWidth = 3
)

// Default cursor color (medium gray)
var (
	DefaultCursorColor = color.NRGBA{R: 80, G: 80, B: 80, A: 255}
)

// Cursor renders the brush cursor at the specified canvas coordinates.
// Returns a slice of canvas objects representing the cursor visualization.
func Cursor(config apptype.PelCanvasConfig, brushType apptype.BrushType, ev *desktop.MouseEvent, x int, y int) []fyne.CanvasObject {
	var objects []fyne.CanvasObject

	switch brushType {
	case apptype.BrushTypePencil, apptype.BrushTypeEraser:
		objects = renderPixelCursor(config, x, y)
	case apptype.BrushTypeFill:
		objects = renderFillCursor(config, x, y)
	case apptype.BrushTypeLine:
		objects = renderLineCursor(config, x, y)
	case apptype.BrushTypeRectangle:
		objects = renderRectangleCursor(config, x, y)
	case apptype.BrushTypeCircle:
		objects = renderCircleCursor(config, x, y)
	default:
		objects = renderPixelCursor(config, x, y)
	}

	return objects
}

// renderPixelCursor creates a square cursor outline for pixel-based tools
func renderPixelCursor(config apptype.PelCanvasConfig, x, y int) []fyne.CanvasObject {
	pxSize := float32(config.PxSize)
	xOrigin := (float32(x) * pxSize) + config.CanvasOffset.X
	yOrigin := (float32(y) * pxSize) + config.CanvasOffset.Y

	// Create four lines forming a square
	lines := make([]fyne.CanvasObject, 4)

	// Left edge
	lines[0] = createCursorLine(
		fyne.NewPos(xOrigin, yOrigin),
		fyne.NewPos(xOrigin, yOrigin+pxSize),
	)

	// Top edge
	lines[1] = createCursorLine(
		fyne.NewPos(xOrigin, yOrigin),
		fyne.NewPos(xOrigin+pxSize, yOrigin),
	)

	// Right edge
	lines[2] = createCursorLine(
		fyne.NewPos(xOrigin+pxSize, yOrigin),
		fyne.NewPos(xOrigin+pxSize, yOrigin+pxSize),
	)

	// Bottom edge
	lines[3] = createCursorLine(
		fyne.NewPos(xOrigin, yOrigin+pxSize),
		fyne.NewPos(xOrigin+pxSize, yOrigin+pxSize),
	)

	return lines
}

// renderFillCursor creates a cursor for the fill tool
func renderFillCursor(config apptype.PelCanvasConfig, x, y int) []fyne.CanvasObject {
	// For now, use pixel cursor
	// TODO: Implement custom fill cursor (e.g., paint bucket icon)
	return renderPixelCursor(config, x, y)
}

// renderLineCursor creates a cursor for the line tool
func renderLineCursor(config apptype.PelCanvasConfig, x, y int) []fyne.CanvasObject {
	// For now, use pixel cursor
	// TODO: Implement line tool with start/end point visualization
	return renderPixelCursor(config, x, y)
}

// renderRectangleCursor creates a cursor for the rectangle tool
func renderRectangleCursor(config apptype.PelCanvasConfig, x, y int) []fyne.CanvasObject {
	// For now, use pixel cursor
	// TODO: Implement rectangle preview
	return renderPixelCursor(config, x, y)
}

// renderCircleCursor creates a cursor for the circle tool
func renderCircleCursor(config apptype.PelCanvasConfig, x, y int) []fyne.CanvasObject {
	// For now, use pixel cursor
	// TODO: Implement circle preview
	return renderPixelCursor(config, x, y)
}

// createCursorLine creates a styled line for cursor rendering
func createCursorLine(pos1, pos2 fyne.Position) *canvas.Line {
	line := canvas.NewLine(DefaultCursorColor)
	line.StrokeWidth = DefaultCursorStrokeWidth
	line.Position1 = pos1
	line.Position2 = pos2
	return line
}

// TryBrush attempts to apply the current brush tool based on mouse events.
// Returns true if the brush operation was successful, false otherwise.
func TryBrush(appState *apptype.State, brushable apptype.Brushable, ev *desktop.MouseEvent) bool {
	if appState == nil || brushable == nil || ev == nil {
		return false
	}

	switch appState.BrushType {
	case apptype.BrushTypePencil:
		return tryPaintPixel(appState, brushable, ev)
	case apptype.BrushTypeEraser:
		return tryErasePixel(appState, brushable, ev)
	case apptype.BrushTypeFill:
		return tryFillArea(appState, brushable, ev)
	case apptype.BrushTypeLine:
		return tryDrawLine(appState, brushable, ev)
	case apptype.BrushTypeRectangle:
		return tryDrawRectangle(appState, brushable, ev)
	case apptype.BrushTypeCircle:
		return tryDrawCircle(appState, brushable, ev)
	default:
		return false
	}
}

// tryPaintPixel paints a single pixel with the current brush color
func tryPaintPixel(appState *apptype.State, brushable apptype.Brushable, ev *desktop.MouseEvent) bool {
	x, y := brushable.MouseToCanvasXY(ev)
	if x != nil && y != nil && ev.Button == desktop.MouseButtonPrimary {
		if err := brushable.SetColor(appState.BrushColor, *x, *y); err != nil {
			// Log error but don't fail - just return false
			return false
		}
		return true
	}
	return false
}

// tryErasePixel erases a pixel by setting it to transparent
func tryErasePixel(appState *apptype.State, brushable apptype.Brushable, ev *desktop.MouseEvent) bool {
	x, y := brushable.MouseToCanvasXY(ev)
	if x != nil && y != nil && ev.Button == desktop.MouseButtonPrimary {
		transparentColor := color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		if err := brushable.SetColor(transparentColor, *x, *y); err != nil {
			return false
		}
		return true
	}
	return false
}

// tryFillArea fills an area with the current brush color (flood fill)
func tryFillArea(appState *apptype.State, brushable apptype.Brushable, ev *desktop.MouseEvent) bool {
	// TODO: Implement flood fill algorithm
	x, y := brushable.MouseToCanvasXY(ev)
	if x != nil && y != nil && ev.Button == desktop.MouseButtonPrimary {
		// Placeholder: just paint single pixel for now
		if err := brushable.SetColor(appState.BrushColor, *x, *y); err != nil {
			return false
		}
		return true
	}
	return false
}

// tryDrawLine draws a line between two points
func tryDrawLine(appState *apptype.State, brushable apptype.Brushable, ev *desktop.MouseEvent) bool {
	// TODO: Implement line drawing with start/end point tracking
	return tryPaintPixel(appState, brushable, ev)
}

// tryDrawRectangle draws a rectangle
func tryDrawRectangle(appState *apptype.State, brushable apptype.Brushable, ev *desktop.MouseEvent) bool {
	// TODO: Implement rectangle drawing
	return tryPaintPixel(appState, brushable, ev)
}

// tryDrawCircle draws a circle
func tryDrawCircle(appState *apptype.State, brushable apptype.Brushable, ev *desktop.MouseEvent) bool {
	// TODO: Implement circle drawing
	return tryPaintPixel(appState, brushable, ev)
}
