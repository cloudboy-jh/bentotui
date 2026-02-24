package theme

type Theme struct {
	Accent     string
	Text       string
	Muted      string
	Background string
	Success    string
	Warning    string
	Error      string

	Surface       string
	SurfaceMuted  string
	Border        string
	BorderFocused string
	TitleText     string
	TitleBG       string
	StatusText    string
	StatusBG      string
	DialogText    string
	DialogBG      string
	DialogBorder  string
	Scrim         string
}

type Option func(*Theme)

func New(opts ...Option) Theme {
	t := Preset("amber")
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

func Preset(name string) Theme {
	switch name {
	case "emerald":
		return Theme{
			Accent:        "#10B981",
			Text:          "#ECFDF5",
			Muted:         "#6EE7B7",
			Background:    "#052E1C",
			Success:       "#22C55E",
			Warning:       "#F59E0B",
			Error:         "#EF4444",
			Surface:       "#083122",
			SurfaceMuted:  "#0D3E2C",
			Border:        "#1D6A4E",
			BorderFocused: "#34D399",
			TitleText:     "#ECFDF5",
			TitleBG:       "#0F513A",
			StatusText:    "#ECFDF5",
			StatusBG:      "#0F513A",
			DialogText:    "#ECFDF5",
			DialogBG:      "#083122",
			DialogBorder:  "#34D399",
			Scrim:         "#031A11",
		}
	case "violet":
		return Theme{
			Accent:        "#8B5CF6",
			Text:          "#F5F3FF",
			Muted:         "#C4B5FD",
			Background:    "#1E1B4B",
			Success:       "#10B981",
			Warning:       "#F59E0B",
			Error:         "#EF4444",
			Surface:       "#241F61",
			SurfaceMuted:  "#2D2774",
			Border:        "#4C3FA4",
			BorderFocused: "#A78BFA",
			TitleText:     "#F5F3FF",
			TitleBG:       "#372E84",
			StatusText:    "#F5F3FF",
			StatusBG:      "#372E84",
			DialogText:    "#F5F3FF",
			DialogBG:      "#241F61",
			DialogBorder:  "#A78BFA",
			Scrim:         "#151235",
		}
	default:
		return Theme{
			Accent:        "#D9A35B",
			Text:          "#E6EDF3",
			Muted:         "#8B97A8",
			Background:    "#0D1117",
			Success:       "#10B981",
			Warning:       "#F59E0B",
			Error:         "#EF4444",
			Surface:       "#161B22",
			SurfaceMuted:  "#1C2430",
			Border:        "#2D3642",
			BorderFocused: "#D9A35B",
			TitleText:     "#0D1117",
			TitleBG:       "#D9A35B",
			StatusText:    "#E6EDF3",
			StatusBG:      "#1B202A",
			DialogText:    "#E6EDF3",
			DialogBG:      "#161B22",
			DialogBorder:  "#D9A35B",
			Scrim:         "#0A0E14",
		}
	}
}

func Accent(v string) Option     { return func(t *Theme) { t.Accent = v } }
func Text(v string) Option       { return func(t *Theme) { t.Text = v } }
func Muted(v string) Option      { return func(t *Theme) { t.Muted = v } }
func Background(v string) Option { return func(t *Theme) { t.Background = v } }
func Success(v string) Option    { return func(t *Theme) { t.Success = v } }
func Warning(v string) Option    { return func(t *Theme) { t.Warning = v } }
func Error(v string) Option      { return func(t *Theme) { t.Error = v } }
