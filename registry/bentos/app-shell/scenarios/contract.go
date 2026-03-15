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
			ID:          "hierarchy",
			Title:       "Hierarchy",
			Description: "panel depth and focus separation with strict row paint ownership",
			Run:         runHierarchy,
		},
		{
			ID:          "footer",
			Title:       "Footer",
			Description: "anchored card readability and overflow under constrained width",
			Run:         runFooter,
		},
		{
			ID:          "list",
			Title:       "List",
			Description: "status-heavy rows with ANSI truncation and right-stat alignment",
			Run:         runList,
		},
		{
			ID:          "overlay",
			Title:       "Overlay",
			Description: "dialog layering over busy surfaces with stable footer anchoring",
			Run:         runOverlay,
		},
		{
			ID:          "stress",
			Title:       "Stress",
			Description: "resize/theme/focus churn checks for seam and clipping regressions",
			Run:         runStress,
		},
	}
}
