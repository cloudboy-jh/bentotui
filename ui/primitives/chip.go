package primitives

import (
	"strings"

	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

func ActionChip(t theme.Theme, variant string, enabled bool, key, label string, keyOnly bool) string {
	sys := styles.New(t)
	keyPart := sys.FooterActionKey(variant, enabled).Render(key)
	if keyOnly || strings.TrimSpace(label) == "" {
		return keyPart
	}
	return keyPart + " " + sys.FooterActionLabel(variant, enabled).Render(label)
}
