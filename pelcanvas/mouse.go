package pelcanvas

import (
	"pel/pelcanvas/brush"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func (pelCanvas *PelCanvas) Scrolled(ev *fyne.ScrollEvent) {
	pelCanvas.scale(int(ev.Scrolled.DY))
	pelCanvas.Refresh()
}

func (pelCanvas *PelCanvas) MouseMoved(ev *desktop.MouseEvent) {
	if x, y := pelCanvas.MouseToCanvasXY(ev); x != nil && y != nil {
		brush.TryBrush(pelCanvas.appState, pelCanvas, ev)
		cursor := brush.Cursor(pelCanvas.PelCanvasConfig, pelCanvas.appState.BrushType, ev, *x, *y)
		pelCanvas.renderer.SetCursor(cursor)
	} else {
		pelCanvas.renderer.SetCursor(make([]fyne.CanvasObject, 0))
	}
	pelCanvas.TryPan(pelCanvas.mouseState.previousCoord, ev)
	pelCanvas.Refresh()
	pelCanvas.mouseState.previousCoord = &ev.PointEvent
}

func (pelCanvas *PelCanvas) MouseDown(ev *desktop.MouseEvent) {
	brush.TryBrush(pelCanvas.appState, pelCanvas, ev)
}

func (pelCanvas *PelCanvas) MouseIn(ev *desktop.MouseEvent) {}
func (pelCanvas *PelCanvas) MouseUp(ev *desktop.MouseEvent) {}
func (pelCanvas *PelCanvas) MouseOut()                      {}
