package state

import (
	tea "charm.land/bubbletea/v2"
	elevatedcard "github.com/cloudboy-jh/bentotui/registry/bricks/elevated-card"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
)

type centerDeck struct {
	width  int
	height int

	primary  *elevatedcard.Model
	checks   *elevatedcard.Model
	metrics  *elevatedcard.Model
	primaryT *textBlock
	checksT  *textBlock
	metricsT *textBlock
}

func newCenterDeck() *centerDeck {
	p := &textBlock{}
	c := &textBlock{}
	m := &textBlock{}
	return &centerDeck{
		primary:  elevatedcard.New(elevatedcard.Title("Scenario Output"), elevatedcard.Content(p), elevatedcard.Inset(1)),
		checks:   elevatedcard.New(elevatedcard.Title("Checks"), elevatedcard.Content(c), elevatedcard.Inset(1)),
		metrics:  elevatedcard.New(elevatedcard.Title("Context"), elevatedcard.Content(m), elevatedcard.Inset(1)),
		primaryT: p,
		checksT:  c,
		metricsT: m,
	}
}

func (d *centerDeck) Init() tea.Cmd                           { return nil }
func (d *centerDeck) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return d, nil }

func (d *centerDeck) View() tea.View {
	if d.width <= 0 || d.height <= 0 {
		return tea.NewView("")
	}
	topH := max(3, (d.height*2)/3)
	if topH >= d.height {
		topH = d.height - 1
	}
	if topH < 1 {
		topH = 1
	}
	bottomH := max(1, d.height-topH)

	bottom := rooms.HSplit(d.width, bottomH, d.checks, d.metrics, rooms.WithGutter(1))
	return tea.NewView(rooms.BigTopStrip(d.width, d.height, bottomH, d.primary, rooms.Static(bottom)))
}

func (d *centerDeck) SetSize(width, height int) {
	d.width = width
	d.height = height
	topH := max(3, (height*2)/3)
	if topH >= height {
		topH = height - 1
	}
	if topH < 1 {
		topH = 1
	}
	bottomH := max(1, height-topH)
	d.primary.SetSize(width, topH)
	leftW := max(1, (width-1)/2)
	rightW := max(1, width-leftW-1)
	d.checks.SetSize(leftW, bottomH)
	d.metrics.SetSize(rightW, bottomH)
}

func (d *centerDeck) SetOutput(title, body string) {
	d.primary.SetTitle(title)
	d.primaryT.SetText(body)
}

func (d *centerDeck) SetChecks(v string) {
	d.checksT.SetText(v)
}

func (d *centerDeck) SetMetrics(v string) {
	d.metricsT.SetText(v)
}
