package theme

import "fmt"

type SurfaceTokens struct {
	Canvas      string
	Panel       string
	Elevated    string
	Overlay     string
	Interactive string
}

type TextTokens struct {
	Primary string
	Muted   string
	Inverse string
	Accent  string
}

type BorderTokens struct {
	Normal string
	Subtle string
	Focus  string
}

type StateTokens struct {
	Info    string
	Success string
	Warning string
	Danger  string
}

type SelectionTokens struct {
	BG string
	FG string
}

type InputTokens struct {
	BG          string
	FG          string
	Placeholder string
	Cursor      string
	Border      string
}

type BarTokens struct {
	BG string
	FG string
}

type DialogTokens struct {
	BG     string
	FG     string
	Border string
	Scrim  string
}

type Theme struct {
	Name      string
	Surface   SurfaceTokens
	Text      TextTokens
	Border    BorderTokens
	State     StateTokens
	Selection SelectionTokens
	Input     InputTokens
	Bar       BarTokens
	Dialog    DialogTokens
}

const (
	DefaultName         = "catppuccin-mocha"
	CatppuccinMochaName = "catppuccin-mocha"
	DraculaName         = "dracula"
	OsakaJadeName       = "osaka-jade"
)

func AvailableThemes() []string {
	return availableThemeNames()
}

func Preset(name string) Theme {
	t, ok := presetTheme(name)
	if !ok {
		t, _ = presetTheme(DefaultName)
	}
	return t
}

func presetTheme(name string) (Theme, bool) {
	t, ok := builtinThemes[name]
	if !ok {
		return Theme{}, false
	}
	if t.Name == "" {
		t.Name = name
	}
	return t, true
}

func validateTheme(t Theme) error {
	checks := []struct {
		label string
		value string
	}{
		{"surface.canvas", t.Surface.Canvas},
		{"surface.panel", t.Surface.Panel},
		{"surface.elevated", t.Surface.Elevated},
		{"surface.overlay", t.Surface.Overlay},
		{"surface.interactive", t.Surface.Interactive},
		{"text.primary", t.Text.Primary},
		{"text.muted", t.Text.Muted},
		{"text.inverse", t.Text.Inverse},
		{"text.accent", t.Text.Accent},
		{"border.normal", t.Border.Normal},
		{"border.subtle", t.Border.Subtle},
		{"border.focus", t.Border.Focus},
		{"state.info", t.State.Info},
		{"state.success", t.State.Success},
		{"state.warning", t.State.Warning},
		{"state.danger", t.State.Danger},
		{"selection.bg", t.Selection.BG},
		{"selection.fg", t.Selection.FG},
		{"input.bg", t.Input.BG},
		{"input.fg", t.Input.FG},
		{"input.placeholder", t.Input.Placeholder},
		{"input.cursor", t.Input.Cursor},
		{"input.border", t.Input.Border},
		{"bar.bg", t.Bar.BG},
		{"bar.fg", t.Bar.FG},
		{"dialog.bg", t.Dialog.BG},
		{"dialog.fg", t.Dialog.FG},
		{"dialog.border", t.Dialog.Border},
		{"dialog.scrim", t.Dialog.Scrim},
	}
	for _, c := range checks {
		if c.value == "" {
			return fmt.Errorf("theme token %q is required", c.label)
		}
	}
	return nil
}

var builtinThemes = map[string]Theme{
	CatppuccinMochaName: {
		Name: CatppuccinMochaName,
		Surface: SurfaceTokens{
			Canvas:      "#181825",
			Panel:       "#24273A",
			Elevated:    "#313244",
			Overlay:     "#1E1E2E",
			Interactive: "#2B2C3F",
		},
		Text: TextTokens{
			Primary: "#CDD6F4",
			Muted:   "#BAC2DE",
			Inverse: "#1E1E2E",
			Accent:  "#89B4FA",
		},
		Border: BorderTokens{
			Normal: "#585B70",
			Subtle: "#45475A",
			Focus:  "#89B4FA",
		},
		State: StateTokens{
			Info:    "#89B4FA",
			Success: "#A6E3A1",
			Warning: "#F9E2AF",
			Danger:  "#F38BA8",
		},
		Selection: SelectionTokens{BG: "#89B4FA", FG: "#1E1E2E"},
		Input: InputTokens{
			BG:          "#2B2C3F",
			FG:          "#CDD6F4",
			Placeholder: "#BAC2DE",
			Cursor:      "#89B4FA",
			Border:      "#74C7EC",
		},
		Bar:    BarTokens{BG: "#11111B", FG: "#CDD6F4"},
		Dialog: DialogTokens{BG: "#313244", FG: "#CDD6F4", Border: "#89B4FA", Scrim: "#0F0F17"},
	},
	DraculaName: {
		Name: DraculaName,
		Surface: SurfaceTokens{
			Canvas:      "#282A36",
			Panel:       "#303341",
			Elevated:    "#343746",
			Overlay:     "#2C3040",
			Interactive: "#3A3D4C",
		},
		Text: TextTokens{
			Primary: "#F8F8F2",
			Muted:   "#B2BEDC",
			Inverse: "#1E1F29",
			Accent:  "#FF79C6",
		},
		Border:    BorderTokens{Normal: "#6272A4", Subtle: "#4F5C88", Focus: "#FF79C6"},
		State:     StateTokens{Info: "#BD93F9", Success: "#50FA7B", Warning: "#FFB86C", Danger: "#FF5555"},
		Selection: SelectionTokens{BG: "#BD93F9", FG: "#1E1F29"},
		Input: InputTokens{
			BG:          "#3A3D4C",
			FG:          "#F8F8F2",
			Placeholder: "#B2BEDC",
			Cursor:      "#FF79C6",
			Border:      "#BD93F9",
		},
		Bar:    BarTokens{BG: "#1F2230", FG: "#F8F8F2"},
		Dialog: DialogTokens{BG: "#2F3343", FG: "#F8F8F2", Border: "#FF79C6", Scrim: "#161821"},
	},
	OsakaJadeName: {
		Name: OsakaJadeName,
		Surface: SurfaceTokens{
			Canvas:      "#071B1A",
			Panel:       "#0C2322",
			Elevated:    "#0D2A28",
			Overlay:     "#0F2927",
			Interactive: "#13322E",
		},
		Text: TextTokens{
			Primary: "#D5EFE9",
			Muted:   "#86B8AC",
			Inverse: "#071B1A",
			Accent:  "#38C2A3",
		},
		Border:    BorderTokens{Normal: "#2F6E63", Subtle: "#244F48", Focus: "#5DE0BF"},
		State:     StateTokens{Info: "#38C2A3", Success: "#56D39B", Warning: "#F4C16D", Danger: "#F26A6A"},
		Selection: SelectionTokens{BG: "#1B5049", FG: "#E3FBF5"},
		Input: InputTokens{
			BG:          "#13322E",
			FG:          "#D5EFE9",
			Placeholder: "#86B8AC",
			Cursor:      "#5DE0BF",
			Border:      "#3E8C80",
		},
		Bar:    BarTokens{BG: "#0B2422", FG: "#D5EFE9"},
		Dialog: DialogTokens{BG: "#102F2B", FG: "#D5EFE9", Border: "#5DE0BF", Scrim: "#031211"},
	},
}
