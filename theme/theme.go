package theme

type Theme struct {
	Accent     string
	Text       string
	Muted      string
	Background string
	PanelBG    string
	ElementBG  string
	ModalBG    string
	Success    string
	Warning    string
	Error      string

	Surface       string
	SurfaceMuted  string
	Border        string
	BorderSubtle  string
	BorderFocused string
	SelectionBG   string
	SelectionText string
	InputBG       string
	InputBorder   string
	TitleText     string
	TitleBG       string
	StatusText    string
	StatusBG      string
	DialogText    string
	DialogBG      string
	DialogBorder  string
	Scrim         string
}

const (
	DefaultName         = "catppuccin-mocha"
	CatppuccinMochaName = "catppuccin-mocha"
	DraculaName         = "dracula"
	OsakaJadeName       = "osaka-jade"
)

type Option func(*Theme)

func New(opts ...Option) Theme {
	t := Preset(DefaultName)
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

func Preset(name string) Theme {
	switch name {
	case OsakaJadeName:
		return Theme{
			Accent:        "#38C2A3",
			Text:          "#D5EFE9",
			Muted:         "#86B8AC",
			Background:    "#071B1A",
			PanelBG:       "#0C2322",
			ElementBG:     "#13302E",
			ModalBG:       "#0F2927",
			Success:       "#56D39B",
			Warning:       "#F4C16D",
			Error:         "#F26A6A",
			Surface:       "#0D2A28",
			SurfaceMuted:  "#123833",
			Border:        "#2F6E63",
			BorderSubtle:  "#244F48",
			BorderFocused: "#5DE0BF",
			SelectionBG:   "#1B5049",
			SelectionText: "#E3FBF5",
			InputBG:       "#13322E",
			InputBorder:   "#3E8C80",
			TitleText:     "#071B1A",
			TitleBG:       "#38C2A3",
			StatusText:    "#D5EFE9",
			StatusBG:      "#0B2422",
			DialogText:    "#D5EFE9",
			DialogBG:      "#102F2B",
			DialogBorder:  "#5DE0BF",
			Scrim:         "#031211",
		}
	case DraculaName:
		return Theme{
			Accent:        "#FF79C6",
			Text:          "#F8F8F2",
			Muted:         "#B2BEDC",
			Background:    "#282A36",
			PanelBG:       "#303341",
			ElementBG:     "#3B3E4D",
			ModalBG:       "#2C3040",
			Success:       "#50FA7B",
			Warning:       "#FFB86C",
			Error:         "#FF5555",
			Surface:       "#343746",
			SurfaceMuted:  "#44475A",
			Border:        "#6272A4",
			BorderSubtle:  "#4F5C88",
			BorderFocused: "#FF79C6",
			SelectionBG:   "#BD93F9",
			SelectionText: "#1E1F29",
			InputBG:       "#3A3D4C",
			InputBorder:   "#BD93F9",
			TitleText:     "#282A36",
			TitleBG:       "#BD93F9",
			StatusText:    "#F8F8F2",
			StatusBG:      "#1F2230",
			DialogText:    "#F8F8F2",
			DialogBG:      "#2F3343",
			DialogBorder:  "#FF79C6",
			Scrim:         "#161821",
		}
	default:
		return Theme{
			Accent:        "#89B4FA",
			Text:          "#CDD6F4",
			Muted:         "#BAC2DE",
			Background:    "#181825",
			PanelBG:       "#24273A",
			ElementBG:     "#313244",
			ModalBG:       "#1E1E2E",
			Success:       "#A6E3A1",
			Warning:       "#F9E2AF",
			Error:         "#F38BA8",
			Surface:       "#313244",
			SurfaceMuted:  "#45475A",
			Border:        "#585B70",
			BorderSubtle:  "#45475A",
			BorderFocused: "#89B4FA",
			SelectionBG:   "#89B4FA",
			SelectionText: "#1E1E2E",
			InputBG:       "#2B2C3F",
			InputBorder:   "#74C7EC",
			TitleText:     "#1E1E2E",
			TitleBG:       "#89B4FA",
			StatusText:    "#CDD6F4",
			StatusBG:      "#11111B",
			DialogText:    "#CDD6F4",
			DialogBG:      "#313244",
			DialogBorder:  "#89B4FA",
			Scrim:         "#0F0F17",
		}
	}
}

func AvailableThemes() []string {
	return []string{CatppuccinMochaName, DraculaName, OsakaJadeName}
}

func Accent(v string) Option     { return func(t *Theme) { t.Accent = v } }
func Text(v string) Option       { return func(t *Theme) { t.Text = v } }
func Muted(v string) Option      { return func(t *Theme) { t.Muted = v } }
func Background(v string) Option { return func(t *Theme) { t.Background = v } }
func Success(v string) Option    { return func(t *Theme) { t.Success = v } }
func Warning(v string) Option    { return func(t *Theme) { t.Warning = v } }
func Error(v string) Option      { return func(t *Theme) { t.Error = v } }
