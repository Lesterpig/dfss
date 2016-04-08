package gui

import (
	"math"

	"github.com/visualfc/goqt/ui"
)

const ARROW_T = math.Pi / 6
const ARROW_L = 15

func (w *Window) DrawClients() {
	scene := w.graphics.Scene()
	for i, c := range w.scene.Clients {
		x, y := w.GetClientPosition(i)

		// Add ellipse
		scene.AddEllipseFWithXYWidthHeightPenBrush(x-10, y-10, 20, 20, pen_black, brush_black)

		// Add text
		t := scene.AddSimpleText(c.Name)
		r := t.BoundingRect()
		t.SetX(x - r.Width()/2)
		t.SetY(y + 10)
	}
}

func (w *Window) GetClientPosition(i int) (x, y float64) {
	if i < 0 {
		return w.GetServerPosition(i == -1)
	}

	nbClients := float64(len(w.scene.Clients))
	angle := 2 * math.Pi * float64(i) / nbClients
	return math.Cos(angle) * (w.circleSize / 2), math.Sin(angle) * (w.circleSize / 2)
}

func (w *Window) GetServerPosition(platform bool) (x, y float64) {
	x = w.circleSize/2 + 150
	y = 0
	if !platform {
		x *= -1
	}
	return
}

func (w *Window) DrawServers() {
	scene := w.graphics.Scene()

	ttp := scene.AddPixmap(w.pixmaps["ttp"])
	x, y := w.GetServerPosition(false)
	ttp.SetPosFWithXY(x-32, y-16)
	ttp.SetToolTip("TTP")

	platform := scene.AddPixmap(w.pixmaps["platform"])
	x, y = w.GetServerPosition(true)
	platform.SetPosFWithXY(x, y-16)
	platform.SetToolTip("Platform")
}

func (w *Window) DrawArrow(xa, ya, xb, yb float64, rgb uint32) {
	scene := w.graphics.Scene()

	path := ui.NewPainterPath()
	path.MoveToFWithXY(xa, ya)
	path.LineToFWithXY(xb, yb)

	v := ui.NewVector2DWithXposYpos(xa-xb, ya-yb)
	l := v.Length()

	// from http://math.stackexchange.com/a/1314050
	xc := xb + ARROW_L/l*(v.X()*math.Cos(ARROW_T)+v.Y()*math.Sin(ARROW_T))
	yc := yb + ARROW_L/l*(v.Y()*math.Cos(ARROW_T)-v.X()*math.Sin(ARROW_T))
	xd := xb + ARROW_L/l*(v.X()*math.Cos(ARROW_T)-v.Y()*math.Sin(ARROW_T))
	yd := yb + ARROW_L/l*(v.Y()*math.Cos(ARROW_T)+v.X()*math.Sin(ARROW_T))

	path.LineToFWithXY(xc, yc)
	path.LineToFWithXY(xd, yd)
	path.LineToFWithXY(xb, yb)
	path.SetFillRule(ui.Qt_WindingFill)

	color := ui.NewColorWithRgb(rgb)
	color.SetAlpha(200)

	pen := ui.NewPenWithColor(color)
	pen.SetWidth(3)
	pen.SetJoinStyle(ui.Qt_RoundJoin)

	brush := ui.NewBrush()
	brush.SetColor(color)
	brush.SetStyle(ui.Qt_SolidPattern)

	arrow := scene.AddPathWithPathPenBrush(path, pen, brush)
	w.currentArrows = append(w.currentArrows, arrow)
}

func (w *Window) RemoveArrows() {
	scene := w.graphics.Scene()

	for _, arrow := range w.currentArrows {
		scene.RemoveItem(&arrow.QGraphicsItem)
		defer arrow.Delete()
	}

	w.currentArrows = nil
}
