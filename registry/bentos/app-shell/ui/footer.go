package ui

import "github.com/cloudboy-jh/bentotui/registry/bricks/bar"

func FooterCards() []bar.Card {
	return []bar.Card{
		{Command: "up/down", Label: "section", Variant: bar.CardPrimary, Enabled: true, Priority: 7},
		{Command: "left/right", Label: "queue", Variant: bar.CardNormal, Enabled: true, Priority: 6},
		{Command: "enter", Label: "pulse", Variant: bar.CardNormal, Enabled: true, Priority: 5},
		{Command: "ctrl+k", Label: "palette", Variant: bar.CardNormal, Enabled: true, Priority: 4},
		{Command: "t", Label: "theme", Variant: bar.CardNormal, Enabled: true, Priority: 4},
		{Command: "c", Label: "compact", Variant: bar.CardNormal, Enabled: true, Priority: 3},
		{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
	}
}
