package gui

// This file stores useful colors

import "github.com/visualfc/goqt/ui"

var colors = map[string]uint32{
	"red":   0x00ff0000,
	"green": 0x0000aa00,
	"blue":  0x000000ff,
	"black": 0x00000000,
}

var pen_black = ui.NewPenWithColor(ui.NewColorWithGlobalcolor(ui.Qt_black))
var pen_gray = ui.NewPenWithColor(ui.NewColorWithGlobalcolor(ui.Qt_gray))

var brush_none = ui.NewBrushWithGlobalcolorBrushstyle(ui.Qt_transparent, ui.Qt_SolidPattern)
var brush_black = ui.NewBrushWithGlobalcolorBrushstyle(ui.Qt_black, ui.Qt_SolidPattern)
