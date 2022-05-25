package pelcanvas

import "fyne.io/fyne/v2"

func (pelCanvas *PelCanvas) scale(direction int) {
	switch {
	case direction > 0:
		pelCanvas.PxSize += 1
	case direction < 0:
		if pelCanvas.PxSize > 2 {
			pelCanvas.PxSize -= 1
		}
	default:
		pelCanvas.PxSize = 10
	}
}

func (pelCanvas *PelCanvas) Pan(previousCoord, currentCoord fyne.PointEvent) {
	xDiff := currentCoord.Position.X - previousCoord.Position.X
	yDiff := currentCoord.Position.Y - previousCoord.Position.Y

	pelCanvas.CanvasOffset.X += xDiff
	pelCanvas.CanvasOffset.Y += yDiff
	pelCanvas.Refresh()
}
