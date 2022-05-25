package pelcanvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func (pelCanvas *PelCanvas) Scrolled(ev *fyne.ScrollEvent) {
	pelCanvas.scale(int(ev.Scrolled.DY))
	pelCanvas.Refresh()
}

func (pelCanvas *PelCanvas) MouseMoved(ev *desktop.MouseEvent) {
	pelCanvas.TryPan(pelCanvas.mouseState.previousCoord, ev)
	pelCanvas.Refresh()
	pelCanvas.mouseState.previousCoord = &ev.PointEvent
}

func (pelCanvas *PelCanvas) MouseIn(ev *desktop.MouseEvent) {}
func (pelCanvas *PelCanvas) MouseOut()                      {}
