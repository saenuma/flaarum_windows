package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type halfes struct {
}

func (d *halfes) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for i, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		if i == 0 {
			h += childSize.Height
		}
	}
	return fyne.NewSize(w, h)
}

func (d *halfes) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, containerSize.Height-d.MinSize(objects).Height)
	for _, o := range objects {
		size := o.MinSize()
		newWidth := (containerSize.Width - float32(len(objects)*10)) / float32(len(objects))
		newSize := fyne.NewSize(newWidth, size.Height)
		o.Resize(newSize)
		// o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(newSize.Width+10, 0))
	}
}

type diags struct{}

func (d *diags) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(600, 400)
}

func (d *diags) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for i, o := range objects {
		newHeight := o.MinSize().Height
		if _, ok := o.(*widget.Entry); ok {
			newHeight = containerSize.Height - (2 * objects[i-1].MinSize().Height)
		}
		newSize := fyne.NewSize(containerSize.Width, newHeight)
		o.Resize(newSize)
		o.Move(pos)
		pos = pos.Add(fyne.NewPos(0, newHeight+10))
	}
}
