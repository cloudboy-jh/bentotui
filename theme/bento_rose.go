package theme

func bentoRoseTheme() Theme {
	return Theme{
		Name: "bento-rose",
		Surface: SurfaceTokens{
			Canvas:      "#201a19",
			Panel:       "#2f2826",
			Elevated:    "#3a312f",
			Overlay:     "#181312",
			Interactive: "#4a3a35",
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
			BG:          "#4a3a35",
			FG:          "#f4ddca",
			Placeholder: "#c5a791",
			Cursor:      "#ff92b6",
			Border:      "#b990a1",
		},
		Bar: BarTokens{
			BG: "#3a312f",
			FG: "#f4ddca",
		},
		Dialog: DialogTokens{
			BG:     "#4a3a35",
			FG:     "#f4ddca",
			Border: "#b990a1",
			Scrim:  "#181312",
		},
	}
}
