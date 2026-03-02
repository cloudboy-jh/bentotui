package widgets

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

// Box is the internal base for all widgets.
type box struct {
	width    int
	height   int
	padding  padding
	border   lipgloss.Border
	borderFg color.Color
	bg       color.Color
	fg       color.Color
	content  core.Component
}

type padding struct {
	top    int
	right  int
	bottom int
	left   int
}

func newBox() *box {
	return &box{
		border:   lipgloss.HiddenBorder(),
		borderFg: nil,
		bg:       nil,
		fg:       nil,
	}
}

func (b *box) SetSize(width, height int) {
	b.width = width
	b.height = height
	if b.content != nil {
		if sizeable, ok := b.content.(core.Sizeable); ok {
			innerW := max(0, width-b.padding.left-b.padding.right)
			innerH := max(0, height-b.padding.top-b.padding.bottom)
			if b.border != lipgloss.HiddenBorder() {
				innerW = max(0, innerW-2)
				innerH = max(0, innerH-2)
			}
			sizeable.SetSize(innerW, innerH)
		}
	}
}

func (b *box) GetSize() (int, int) {
	return b.width, b.height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
