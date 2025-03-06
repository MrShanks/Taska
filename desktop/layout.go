package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const sideWidth = 220

type taskaLayout struct {
	top, left, content fyne.CanvasObject
	dividers           [2]fyne.CanvasObject
}

func newTaskaLayout(top, left, content fyne.CanvasObject, dividers [2]fyne.CanvasObject) fyne.Layout {
	return &taskaLayout{top: top, left: left, content: content, dividers: dividers}
}

func (t *taskaLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	topHeight := t.top.MinSize().Height

	t.top.Resize(fyne.NewSize(sideWidth, topHeight))

	t.left.Move(fyne.NewPos(0, topHeight))
	t.left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	t.content.Move(fyne.NewPos(sideWidth, topHeight))
	t.content.Resize(fyne.NewSize(size.Width-sideWidth*2, size.Height-topHeight))

	t.dividers[0].Move(fyne.NewPos(0, topHeight))
	t.dividers[0].Resize(fyne.NewSize(size.Width, theme.SeparatorThicknessSize()))

	t.dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
	t.dividers[1].Resize(fyne.NewSize(theme.SeparatorThicknessSize(), size.Height-topHeight))
}

func (t *taskaLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	borders := fyne.NewSize(
		sideWidth*2,
		t.top.MinSize().Height,
	)

	return borders.AddWidthHeight(200, 200)
}
