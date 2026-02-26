package primitives

import (
	"strings"

	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

func Card(t theme.Theme, variant string, enabled bool, command, label string, commandOnly bool) string {
	sys := styles.New(t)
	commandPart := sys.FooterCardCommand(variant, enabled).Render(command)
	if commandOnly || strings.TrimSpace(label) == "" {
		return commandPart
	}
	return commandPart + " " + sys.FooterCardLabel(variant, enabled).Render(label)
}
