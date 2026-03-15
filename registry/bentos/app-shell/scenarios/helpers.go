package scenarios

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

func fitBlock(text string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	lines := strings.Split(text, "\n")
	out := make([]string, 0, height)
	for i := 0; i < height; i++ {
		line := ""
		if i < len(lines) {
			line = styles.ClipANSI(lines[i], width)
		}
		pad := width - lipgloss.Width(line)
		if pad > 0 {
			line += strings.Repeat(" ", pad)
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func ruler(width int) string {
	if width <= 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(width)
	for i := 0; i < width; i++ {
		if i%10 == 0 {
			b.WriteRune('|')
			continue
		}
		b.WriteRune('·')
	}
	return b.String()
}

func rulerIndex(width int) string {
	if width <= 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i < width; i++ {
		if i%10 == 0 {
			b.WriteString(fmt.Sprintf("%d", i/10))
			continue
		}
		b.WriteRune(' ')
	}
	return styles.ClipANSI(b.String(), width)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
