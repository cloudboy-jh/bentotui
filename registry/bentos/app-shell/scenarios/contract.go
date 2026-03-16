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
			ID:          "layout",
			Title:       "Layout",
			Description: "split/drawer/sidebar framing at 80x24, 100x30, and 140x42",
			Run:         runLayout,
		},
		{
			ID:          "footer",
			Title:       "Footer",
			Description: "anchored card readability and overflow under constrained width",
			Run:         runFooter,
		},
		{
			ID:          "stress",
			Title:       "Stress",
			Description: "resize and theme churn checks for seam and clipping regressions",
			Run:         runStress,
		},
	}
}
