package main

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/taylorskalyo/goreader/parse"
)

type pager struct {
	scrollX int
	scrollY int
	doc     parse.Cellbuf
}

// draw displays a pager's cell buffer in the terminal.
func (p pager) draw() error {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	width, height := termbox.Size()
	var centerOffset int
	for y := 0; y < height; y++ {
		for x := 0; x < p.doc.Width; x++ {
			index := (y+p.scrollY)*p.doc.Width + x
			if index >= len(p.doc.Cells) || index <= 0 {
				continue
			}
			cell := p.doc.Cells[index]
			if width > p.doc.Width {
				centerOffset = (width - p.doc.Width) / 2
			}

			// Calling SetCell with coordinates outside of the terminal viewport
			// results in a no-op.
			termbox.SetCell(x+p.scrollX+centerOffset, y, cell.Ch, cell.Fg, cell.Bg)
		}
	}

	return termbox.Flush()
}

// scrollDown pans the pager's viewport down, without exceeding the underlying
// cell buffer document's boundaries.
func (p *pager) scrollDown() {
	if p.scrollY < p.maxScrollY() {
		p.scrollY++
	}
}

// scrollUp pans the pager's viewport up, without exceeding the underlying cell
// buffer document's boundaries.
func (p *pager) scrollUp() {
	if p.scrollY > 0 {
		p.scrollY--
	}
}

// scrollLeft pans the pager's viewport left, without exceeding the underlying
// cell buffer document's boundaries.
func (p *pager) scrollLeft() {
	if p.scrollX < 0 {
		p.scrollX++
	}
}

// scrollRight pans the pager's viewport right, without exceeding the
// underlying cell buffer document's boundaries.
func (p *pager) scrollRight() {
	if p.scrollX > -p.maxScrollX() {
		p.scrollX--
	}
}

// pageDown pans the pager's viewport down by a full page, without exceeding
// the underlying cell buffer document's boundaries.
func (p *pager) pageDown() bool {
	_, viewHeight := termbox.Size()
	if p.scrollY < p.maxScrollY() {
		p.scrollY += viewHeight
		return true
	}

	return false
}

// pageUp pans the pager's viewport up by a full page, without exceeding the
// underlying cell buffer document's boundaries.
func (p *pager) pageUp() bool {
	_, viewHeight := termbox.Size()
	if p.scrollY > viewHeight {
		p.scrollY -= viewHeight
		return true
	} else if p.scrollY > 0 {
		p.scrollY = 0
		return true
	}

	return false
}

// toTop set's the pager's horizontal and vertical panning distance back to
// zero.
func (p *pager) toTop() {
	p.scrollX = 0
	p.scrollY = 0
}

// toBottom set's the pager's horizontal panning distance back to zero and
// vertical panning distance to the last viewport page.
func (p *pager) toBottom() {
	_, viewHeight := termbox.Size()
	p.scrollX = 0
	p.scrollY = p.pages() * viewHeight
}

// maxScrollX represents the pager's maximum horizontal scroll distance.
func (p pager) maxScrollX() int {
	docWidth, _ := p.size()
	viewWidth, _ := termbox.Size()
	return docWidth - viewWidth
}

// maxScrollY represents the pager's maximum vertical scroll distance.
func (p pager) maxScrollY() int {
	_, docHeight := p.size()
	_, viewHeight := termbox.Size()
	return docHeight - viewHeight
}

// size returns the width and height of the pager's underlying cell buffer
// document.
func (p pager) size() (int, int) {
	height := len(p.doc.Cells) / p.doc.Width
	return p.doc.Width, height
}

// pages returns the number of times the pager's underlying cell buffer
// document can be split into viewport sized pages.
func (p pager) pages() int {
	_, docHeight := p.size()
	_, viewHeight := termbox.Size()
	return docHeight / viewHeight
}
