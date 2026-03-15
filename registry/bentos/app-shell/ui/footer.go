package ui

import "github.com/cloudboy-jh/bentotui/registry/bricks/bar"

func FooterCards() []bar.Card {
	return []bar.Card{
		{Command: "j/k", Label: "scenario", Variant: bar.CardPrimary, Enabled: true, Priority: 9},
		{Command: "h/l", Label: "viewport", Variant: bar.CardNormal, Enabled: true, Priority: 8},
		{Command: "[/]", Label: "focus", Variant: bar.CardNormal, Enabled: true, Priority: 7},
		{Command: "t", Label: "theme", Variant: bar.CardNormal, Enabled: true, Priority: 6},
		{Command: "d", Label: "paint", Variant: bar.CardNormal, Enabled: true, Priority: 5},
		{Command: "s", Label: "snapshot", Variant: bar.CardNormal, Enabled: true, Priority: 4},
		{Command: "m", Label: "keymap", Variant: bar.CardNormal, Enabled: true, Priority: 3},
		{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
	}
}
