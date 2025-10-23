// Package swatch provides mouse event handling for color swatch widgets.
package swatch

import "fyne.io/fyne/v2/driver/desktop"

// MouseDown handles mouse button press events on the swatch
// Only responds to the primary (left) mouse button
func (swatch *Swatch) MouseDown(ev *desktop.MouseEvent) {
	if swatch == nil || ev == nil {
		return
	}

	// Only respond to primary mouse button
	if ev.Button != desktop.MouseButtonPrimary {
		return
	}

	// Invoke click handler if present
	if swatch.clickHandler != nil {
		swatch.clickHandler(swatch)
	}

	// Mark this swatch as selected
	swatch.Selected = true
	swatch.Refresh()
}

// MouseUp handles mouse button release events on the swatch
func (swatch *Swatch) MouseUp(ev *desktop.MouseEvent) {
	// TODO: Implement if needed for drag-and-drop color operations
}

// MouseIn handles mouse entering the swatch area
func (swatch *Swatch) MouseIn(ev *desktop.MouseEvent) {
	// TODO: Implement hover effects (e.g., brighten border, show tooltip)
}

// MouseOut handles mouse leaving the swatch area
func (swatch *Swatch) MouseOut() {
	// TODO: Implement hover effects cleanup
}

// MouseMoved handles mouse movement over the swatch
func (swatch *Swatch) MouseMoved(ev *desktop.MouseEvent) {
	// No action needed for basic swatch
}
