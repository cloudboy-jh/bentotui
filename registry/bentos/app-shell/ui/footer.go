package ui

import "github.com/cloudboy-jh/bentotui/registry/bricks/bar"

func FooterCards() []bar.Card {
	return []bar.Card{
		{Command: "up/down", Label: "scenario", Variant: bar.CardPrimary, Enabled: true, Priority: 6},
		{Command: "left/right", Label: "viewport", Variant: bar.CardNormal, Enabled: true, Priority: 5},
		{Command: "t", Label: "theme", Variant: bar.CardNormal, Enabled: true, Priority: 4},
		{Command: "d", Label: "paint", Variant: bar.CardNormal, Enabled: true, Priority: 3},
		{Command: "s", Label: "snapshot", Variant: bar.CardNormal, Enabled: true, Priority: 2},
		{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
	}
}
