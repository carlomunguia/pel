package pelcanvas

import (
	"image"
	"image/color"
	"pel/apptype"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type PelCanvasMouseState struct {
	previousCoord *fyne.PointEvent
}

type PelCanvas struct {
	widget.BaseWidget
	apptype.PelCanvasConfig
	renderer    *PelCanvasRenderer
	PixelData   image.Image
	mouseState  PelCanvasMouseState
	appState    *apptype.State
	reloadImage bool
}

func (pelCanvas *PelCanvas) Bounds() image.Rectangle {
	x0 := int(pelCanvas.CanvasOffset.X)
	y0 := int(pelCanvas.CanvasOffset.Y)
	x1 := int(pelCanvas.PxCols*pelCanvas.PxSize + int(pelCanvas.CanvasOffset.X))
	y1 := int(pelCanvas.PxRows*pelCanvas.PxSize + int(pelCanvas.CanvasOffset.Y))
	return image.Rect(x0, y0, x1, y1)
}

func InBounds(pos fyne.Position, bounds image.Rectangle) bool {
	if pos.X >= float32(bounds.Min.X) &&
		pos.X < float32(bounds.Max.X) &&
		pos.Y >= float32(bounds.Min.Y) &&
		pos.Y < float32(bounds.Max.Y) {
		return true
	}
	return false
}

func NewBlankImage(cols, rows int, c color.Color) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, cols, rows))
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func NewPelCanvas(state *apptype.State, config apptype.PelCanvasConfig) *PelCanvas {
	pelCanvas := &PelCanvas{
		PelCanvasConfig: config,
		appState:        state,
	}
	pelCanvas.PixelData = NewBlankImage(pelCanvas.PxCols, pelCanvas.PxRows, color.NRGBA{128, 128, 128, 255})
	pelCanvas.ExtendBaseWidget(pelCanvas)
	return pelCanvas
}

func (pelCanvas *PelCanvas) CreateRenderer() fyne.WidgetRenderer {
	canvasImage := canvas.NewImageFromImage(pelCanvas.PixelData)
	canvasImage.ScaleMode = canvas.ImageScalePixels
	canvasImage.FillMode = canvas.ImageFillContain

	canvasBorder := make([]canvas.Line, 4)
	for i := 0; i < len(canvasBorder); i++ {
		canvasBorder[i].StrokeColor = color.NRGBA64{100, 100, 100, 255}
		canvasBorder[i].StrokeWidth = 2
	}
	renderer := &PelCanvasRenderer{
		pelCanvas:    pelCanvas,
		canvasImage:  canvasImage,
		canvasBorder: canvasBorder,
	}
	pelCanvas.renderer = renderer
	return renderer
}

func (pelCanvas *PelCanvas) TryPan(previousCoord *fyne.PointEvent, ev *desktop.MouseEvent) {
	if previousCoord != nil && ev.Button == desktop.MouseButtonTertiary {
		pelCanvas.Pan(*previousCoord, ev.PointEvent)
	}
}

func (pelCanvas *PelCanvas) SetColor(c color.Color, x, y int) {
	if nrgba, ok := pelCanvas.PixelData.(*image.NRGBA); ok {
		nrgba.Set(x, y, c)
	}
	if rgba, ok := pelCanvas.PixelData.(*image.RGBA); ok {
		rgba.Set(x, y, c)
	}
	pelCanvas.Refresh()
}

func (pelCanvas *PelCanvas) MouseToCanvasXY(ev *desktop.MouseEvent) (*int, *int) {
	bounds := pelCanvas.Bounds()
	if !InBounds(ev.Position, bounds) {
		return nil, nil
	}
	pxSize := float32(pelCanvas.PxSize)
	xOffset := pelCanvas.CanvasOffset.X
	yOffset := pelCanvas.CanvasOffset.Y

	x := int((ev.Position.X - xOffset) / pxSize)
	y := int((ev.Position.Y - yOffset) / pxSize)

	return &x, &y
}

func (pelCanvas *PelCanvas) LoadImage(img image.Image) {
	dimensions := img.Bounds()

	pelCanvas.PelCanvasConfig.PxCols = dimensions.Dx()
	pelCanvas.PelCanvasConfig.PxRows = dimensions.Dy()

	pelCanvas.PixelData = img
	pelCanvas.reloadImage = true
	pelCanvas.Refresh()
}

func (pelCanvas *PelCanvas) NewDrawing(cols, rows int) {
	pelCanvas.appState.SetFilePath("")
	pelCanvas.PxCols = cols
	pelCanvas.PxRows = rows
	pixelData := NewBlankImage(cols, rows, color.NRGBA{128, 128, 128, 255})
	pelCanvas.LoadImage(pixelData)
}
