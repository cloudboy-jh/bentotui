package scenarios

type Viewport struct {
	Name   string
	Width  int
	Height int
}

type Context struct {
	Width      int
	Height     int
	Viewport   Viewport
	PaintDebug bool
	Snapshot   bool
	FocusOwner string
	StressStep int
}

type CheckLevel string

const (
	CheckPass CheckLevel = "pass"
	CheckWarn CheckLevel = "warn"
	CheckFail CheckLevel = "fail"
)

type Check struct {
	Name   string
	Level  CheckLevel
	Detail string
}

type Result struct {
	Canvas  string
	Checks  []Check
	Metrics map[string]string
}

type Definition struct {
	ID          string
	Title       string
	Description string
	Run         func(ctx Context) Result
}

func All() []Definition {
	return []Definition{
		{
			ID:          "cards-list",
			Title:       "Cards + List",
			Description: "elevated card hosting a list with sections, selection, and stats",
			Run:         runList,
		},
		{
			ID:          "cards-table",
			Title:       "Cards + Table",
			Description: "elevated card hosting a table with aligned columns",
			Run:         runTable,
		},
		{
			ID:          "cards-modal",
			Title:       "Cards + Modal",
			Description: "modal overlay above card content with stable footer lane",
			Run:         runModal,
		},
		{
			ID:          "cards-footer",
			Title:       "Cards + Footer",
			Description: "anchored command footer readability and truncation behavior",
			Run:         runFooter,
		},
	}
}
