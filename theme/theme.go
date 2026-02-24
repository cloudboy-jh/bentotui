package theme

type Theme struct {
	Accent     string
	Text       string
	Muted      string
	Background string
	Success    string
	Warning    string
	Error      string
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
			Accent:     "#10B981",
			Text:       "#ECFDF5",
			Muted:      "#6EE7B7",
			Background: "#052E1C",
			Success:    "#22C55E",
			Warning:    "#F59E0B",
			Error:      "#EF4444",
		}
	case "violet":
		return Theme{
			Accent:     "#8B5CF6",
			Text:       "#F5F3FF",
			Muted:      "#C4B5FD",
			Background: "#1E1B4B",
			Success:    "#10B981",
			Warning:    "#F59E0B",
			Error:      "#EF4444",
		}
	default:
		return Theme{
			Accent:     "#F59E0B",
			Text:       "#F8FAFC",
			Muted:      "#94A3B8",
			Background: "#171717",
			Success:    "#10B981",
			Warning:    "#F59E0B",
			Error:      "#EF4444",
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
