package engine

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type Sizing int

const (
	Fixed Sizing = iota
	Fill
	Ratio
)

type Spec struct {
	Kind Sizing
	N    int
}

func RenderVertical(width, height int, specs []Spec, cells []Sizable) string {
	if width <= 0 || height <= 0 || len(specs) == 0 || len(specs) != len(cells) {
		return ""
	}

	heights := Allocate(specs, height)
	out := make([]string, len(cells))
	for i := range cells {
		h := Max(1, heights[i])
		cells[i].SetSize(width, h)
		out[i] = Constrain(ViewString(cells[i].View()), width, h)
	}

	return Constrain(lipgloss.JoinVertical(lipgloss.Left, out...), width, height)
}

func RenderHorizontal(width, height int, specs []Spec, cells []Sizable) string {
	if width <= 0 || height <= 0 || len(specs) == 0 || len(specs) != len(cells) {
		return ""
	}

	widths := Allocate(specs, width)
	out := make([]string, len(cells))
	for i := range cells {
		w := Max(1, widths[i])
		cells[i].SetSize(w, height)
		out[i] = Constrain(ViewString(cells[i].View()), w, height)
	}

	return Constrain(lipgloss.JoinHorizontal(lipgloss.Top, out...), width, height)
}

func Allocate(specs []Spec, total int) []int {
	if len(specs) == 0 {
		return nil
	}

	total = Max(1, total)
	out := make([]int, len(specs))

	fixedTotal := 0
	ratioWeight := 0
	fillCount := 0

	for i, s := range specs {
		size := Max(1, s.N)
		switch s.Kind {
		case Fixed:
			out[i] = size
			fixedTotal += size
		case Ratio:
			ratioWeight += size
		case Fill:
			fillCount++
		}
	}

	remaining := total - fixedTotal

	if ratioWeight > 0 {
		ratioAssigned := 0
		lastRatio := -1
		for i, s := range specs {
			if s.Kind != Ratio {
				continue
			}
			lastRatio = i
			size := Max(1, s.N)
			part := (remaining * size) / ratioWeight
			part = Max(1, part)
			out[i] = part
			ratioAssigned += part
		}
		if lastRatio >= 0 {
			out[lastRatio] += remaining - ratioAssigned
		}
	}

	used := 0
	for _, n := range out {
		used += n
	}
	remaining = total - used

	if fillCount > 0 {
		share := 0
		if remaining > 0 {
			share = remaining / fillCount
		}
		share = Max(1, share)

		lastFill := -1
		for i, s := range specs {
			if s.Kind != Fill {
				continue
			}
			lastFill = i
			out[i] = share
		}

		used = 0
		for _, n := range out {
			used += n
		}

		if lastFill >= 0 {
			out[lastFill] += total - used
		}
	}

	for i := range out {
		if out[i] < 1 {
			out[i] = 1
		}
	}

	return rebalance(out, total)
}

func rebalance(values []int, total int) []int {
	sum := 0
	for _, n := range values {
		sum += n
	}

	if sum == total {
		return values
	}

	if sum < total {
		values[len(values)-1] += total - sum
		return values
	}

	diff := sum - total
	for i := len(values) - 1; i >= 0 && diff > 0; i-- {
		if values[i] <= 1 {
			continue
		}
		take := Min(diff, values[i]-1)
		values[i] -= take
		diff -= take
	}

	return values
}

func Constrain(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	rawLines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	lines := make([]string, height)

	for i := 0; i < height; i++ {
		line := ""
		if i < len(rawLines) {
			line = rawLines[i]
		}
		lines[i] = constrainLine(line, width)
	}

	return strings.Join(lines, "\n")
}

func constrainLine(line string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(line) > width {
		line = ansi.Truncate(line, width, "")
	}
	if lipgloss.Width(line) < width {
		line += strings.Repeat(" ", width-lipgloss.Width(line))
	}
	return line
}

func Overlay(bg, fg string, x, y int) string {
	bgLines := strings.Split(strings.ReplaceAll(bg, "\r\n", "\n"), "\n")
	fgLines := strings.Split(strings.ReplaceAll(fg, "\r\n", "\n"), "\n")

	if len(bgLines) == 0 || len(fgLines) == 0 {
		return bg
	}

	bgW := lipgloss.Width(bgLines[0])
	for i, fgLine := range fgLines {
		target := y + i
		if target < 0 || target >= len(bgLines) {
			continue
		}

		fgW := lipgloss.Width(fgLine)
		if fgW <= 0 {
			continue
		}

		start := Max(0, x)
		end := Min(bgW, x+fgW)
		if start >= end {
			continue
		}

		left := ansi.Cut(bgLines[target], 0, start)
		right := ansi.Cut(bgLines[target], end, bgW)
		fgPart := ansi.Cut(fgLine, start-x, end-x)

		bgLines[target] = constrainLine(left+fgPart+right, bgW)
	}

	return strings.Join(bgLines, "\n")
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
