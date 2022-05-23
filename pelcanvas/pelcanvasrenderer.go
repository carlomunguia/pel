package pelcanvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type PelCanvasRenderer struct {
	pelCanvas    *PelCanvas
	canvasImage  *canvas.Image
	canvasBorder []canvas.Line
}

func (renderer *PelCanvasRenderer) MinSize() fyne.Size {
	return renderer.pelCanvas.DrawingArea
}

func (renderer *PelCanvasRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, 5)
	for i := 0; i < len(renderer.canvasBorder); i++ {
		objects = append(objects, &renderer.canvasBorder[i])
	}
	objects = append(objects, renderer.canvasImage)
	return objects
}

func (renderer *PelCanvasRenderer) Destroy() {}

func (renderer *PelCanvasRenderer) Layout(size fyne.Size) {
	renderer.LayoutCanvas(size)
	renderer.LayoutBorder(size)
}

func (renderer *PelCanvasRenderer) LayoutCanvas(size fyne.Size) {
	imgPxWidth := renderer.pelCanvas.PxCols
	imgPxHeight := renderer.pelCanvas.PxRows
	pxSize := renderer.pelCanvas.PxSize
	renderer.canvasImage.Move(fyne.NewPos(renderer.pelCanvas.CanvasOffset.X, renderer.pelCanvas.CanvasOffset.Y))
	renderer.canvasImage.Resize(fyne.NewSize(float32(imgPxWidth*pxSize), float32(imgPxHeight*pxSize)))
}

func (renderer *PelCanvasRenderer) LayoutBorder(size fyne.Size) {
	offset := renderer.pelCanvas.CanvasOffset
	imgHeight := renderer.canvasImage.Size().Height
	imgWidth := renderer.canvasImage.Size().Width

	left := &renderer.canvasBorder[0]
	left.Position1 = fyne.NewPos(offset.X, offset.Y)
	left.Position2 = fyne.NewPos(offset.X, offset.Y+imgHeight)

	top := &renderer.canvasBorder[1]
	top.Position1 = fyne.NewPos(offset.X, offset.Y)
	top.Position2 = fyne.NewPos(offset.X+imgWidth, offset.Y)

	right := &renderer.canvasBorder[2]
	right.Position1 = fyne.NewPos(offset.X+imgWidth, offset.Y)
	right.Position2 = fyne.NewPos(offset.X+imgWidth, offset.Y+imgHeight)

	bottom := &renderer.canvasBorder[3]
	bottom.Position1 = fyne.NewPos(offset.X, offset.Y+imgHeight)
	bottom.Position2 = fyne.NewPos(offset.X+imgWidth, offset.Y+imgHeight)

}

func (renderer *PelCanvasRenderer) Refresh() {
	if renderer.pelCanvas.reloadImage {
		renderer.canvasImage = canvas.NewImageFromImage(renderer.pelCanvas.PixelData)
		renderer.canvasImage.ScaleMode = canvas.ImageScalePixels
		renderer.canvasImage.FillMode = canvas.ImageFillContain
		renderer.pelCanvas.reloadImage = false
	}
	renderer.Layout(renderer.pelCanvas.Size())
	canvas.Refresh(renderer.canvasImage)
}
