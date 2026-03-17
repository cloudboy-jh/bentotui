package theme

func bentoRoseTheme() Theme {
	return Theme{
		Name: "bento-rose",
		Surface: SurfaceTokens{
			Canvas:      "#201a19",
			Panel:       "#5f4843",
			Overlay:     "#181312",
			Interactive: "#7a5d56",
		},
		Text: TextTokens{
			Primary: "#f4ddca",
			Muted:   "#d5b8a2",
			Inverse: "#201a19",
			Accent:  "#ff92b6",
		},
		Border: BorderTokens{
			Normal: "#b990a1",
			Subtle: "#a77b8f",
			Focus:  "#ff92b6",
		},
		State: StateTokens{
			Info:    "#86b6ff",
			Success: "#9cd67a",
			Warning: "#f1c56b",
			Danger:  "#ff7e9f",
		},
		Selection: SelectionTokens{
			BG: "#ff92b6",
			FG: "#201a19",
		},
		Input: InputTokens{
			BG:          "#7a5d56",
			FG:          "#f4ddca",
			Placeholder: "#c5a791",
			Cursor:      "#ff92b6",
			Border:      "#b990a1",
		},
		Bar: BarTokens{
			BG: "#5f4843",
			FG: "#f4ddca",
		},
		Footer: FooterTokens{
			AnchoredBG:    "#2f2422",
			AnchoredFG:    "#f4ddca",
			AnchoredMuted: "#d5b8a2",
		},
		Dialog: DialogTokens{
			BG:     "#7a5d56",
			FG:     "#f4ddca",
			Border: "#b990a1",
			Scrim:  "#181312",
		},
		Card: CardTokens{
			ChromeBG:    "#3d2d2a",
			BodyBG:      "#5f4843",
			FrameFG:     "#f4ddca",
			FocusEdgeBG: "#ff92b6",
		},
	}
}
