package state

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	elevatedcard "github.com/cloudboy-jh/bentotui/registry/bricks/elevated-card"
	"github.com/cloudboy-jh/bentotui/registry/bricks/list"
	"github.com/cloudboy-jh/bentotui/registry/bricks/progress"
	"github.com/cloudboy-jh/bentotui/registry/bricks/table"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
)

type centerDeck struct {
	width  int
	height int

	tableCard    *elevatedcard.Model
	queueCard    *elevatedcard.Model
	progressCard *elevatedcard.Model

	services *table.Model
	queue    *list.Model
	progress *progressPane
}

func newCenterDeck() *centerDeck {
	t := table.New("SERVICE", "OWNER", "P95", "ERR%", "DEPLOY")
	t.SetCompact(true)
	t.SetBorderless(true)
	t.SetColumnAlign(2, table.AlignRight)
	t.SetColumnAlign(3, table.AlignRight)
	t.SetColumnMinWidth(0, 10)
	t.SetColumnMinWidth(1, 8)
	t.SetColumnMinWidth(4, 8)
	t.SetColumnPriority(4, 5)
	t.SetColumnPriority(1, 4)
	t.AddRow("checkout-api", "kai", "38ms", "0.1", "2m ago")
	t.AddRow("billing-jobs", "jules", "55ms", "0.0", "9m ago")
	t.AddRow("customer-sync", "rani", "112ms", "1.7", "18m ago")
	t.AddRow("event-router", "mina", "47ms", "0.2", "27m ago")

	q := list.New(32)
	q.SetDensity(list.DensityCompact)
	q.AppendSection("ACTIVE")
	q.AppendRow(list.Row{Primary: "Sync configs", Secondary: "running", Tone: list.ToneInfo, RightStat: "now"})
	q.AppendRow(list.Row{Primary: "Warm caches", Secondary: "queued", Tone: list.ToneNeutral, RightStat: "2m"})
	q.AppendSection("BACKLOG")
	q.AppendRow(list.Row{Primary: "Export snapshots", Secondary: "waiting", Tone: list.ToneWarn, RightStat: "9m"})
	q.AppendRow(list.Row{Primary: "Rotate logs", Secondary: "ready", Tone: list.ToneSuccess, RightStat: "13m"})
	q.SetCursor(0)

	p := newProgressPane()

	return &centerDeck{
		tableCard:    elevatedcard.New(elevatedcard.Title("Services"), elevatedcard.Meta("stable view of key metrics"), elevatedcard.Content(t), elevatedcard.CardVariant(elevatedcard.VariantEmphasis), elevatedcard.Inset(1)),
		queueCard:    elevatedcard.New(elevatedcard.Title("Queue"), elevatedcard.Meta("lightweight task lane"), elevatedcard.Content(q), elevatedcard.Inset(1)),
		progressCard: elevatedcard.New(elevatedcard.Title("Progress"), elevatedcard.Meta("simple throughput snapshot"), elevatedcard.Content(p), elevatedcard.Inset(1)),
		services:     t,
		queue:        q,
		progress:     p,
	}
}

func (d *centerDeck) Init() tea.Cmd                           { return nil }
func (d *centerDeck) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return d, nil }

func (d *centerDeck) View() tea.View {
	if d.width <= 0 || d.height <= 0 {
		return tea.NewView("")
	}
	topH := max(6, (d.height*3)/5)
	if topH >= d.height {
		topH = d.height - 1
	}
	if topH < 1 {
		topH = 1
	}
	bottomH := max(1, d.height-topH)

	bottom := rooms.HSplit(d.width, bottomH, d.queueCard, d.progressCard, rooms.WithGutter(1))
	return tea.NewView(rooms.BigTopStrip(d.width, d.height, bottomH, d.tableCard, rooms.Static(bottom)))
}

func (d *centerDeck) SetSize(width, height int) {
	d.width = width
	d.height = height
	topH := max(6, (height*3)/5)
	if topH >= height {
		topH = height - 1
	}
	if topH < 1 {
		topH = 1
	}
	bottomH := max(1, height-topH)

	d.tableCard.SetSize(width, topH)
	leftW := max(1, (width-1)/2)
	rightW := max(1, width-leftW-1)
	d.queueCard.SetSize(leftW, bottomH)
	d.progressCard.SetSize(rightW, bottomH)

	d.services.SetSize(max(24, width-2), max(4, topH-4))
	d.queue.SetSize(max(18, leftW-2), max(4, bottomH-4))
	d.progress.SetSize(max(18, rightW-2), max(4, bottomH-4))
}

func (d *centerDeck) SetActiveSection(section string) {
	d.tableCard.SetVariant(elevatedcard.VariantDefault)
	d.queueCard.SetVariant(elevatedcard.VariantDefault)
	d.progressCard.SetVariant(elevatedcard.VariantDefault)

	switch section {
	case "Services":
		d.tableCard.SetVariant(elevatedcard.VariantEmphasis)
	case "Queue":
		d.queueCard.SetVariant(elevatedcard.VariantEmphasis)
	case "Progress":
		d.progressCard.SetVariant(elevatedcard.VariantEmphasis)
	default:
		d.tableCard.SetVariant(elevatedcard.VariantEmphasis)
	}
}

func (d *centerDeck) SetQueueCursor(i int) {
	d.queue.SetCursor(i)
}

func (d *centerDeck) SetCompact(v bool) {
	d.services.SetCompact(v)
	if v {
		d.tableCard.SetMeta("compact metrics view")
		return
	}
	d.tableCard.SetMeta("stable view of key metrics")
}

func (d *centerDeck) SetWorkspaceMeta(section string, compact bool) {
	density := "comfortable"
	if compact {
		density = "compact"
	}
	d.queueCard.SetFooter(fmt.Sprintf("section %s | %s", section, density))
}

func (d *centerDeck) SetProgress(value float64, label string) {
	d.progress.SetValue(value, label)
	d.progressCard.SetFooter(fmt.Sprintf("overall %.0f%%", value*100))
}

type progressPane struct {
	bar    *progress.Model
	width  int
	height int
	label  string
}

func newProgressPane() *progressPane {
	b := progress.New(24)
	b.SetLabel("sync")
	b.SetValue(0.62)
	return &progressPane{bar: b, label: "2/3 lanes healthy"}
}

func (p *progressPane) Init() tea.Cmd                           { return nil }
func (p *progressPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return p, nil }

func (p *progressPane) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.bar.SetSize(max(10, width-4), 1)
}

func (p *progressPane) SetValue(v float64, label string) {
	p.bar.SetValue(v)
	if strings.TrimSpace(label) != "" {
		p.label = label
	}
}

func (p *progressPane) View() tea.View {
	if p.height <= 0 {
		return tea.NewView("")
	}
	line := ansi.Strip(viewString(p.bar.View()))
	rows := []string{line, p.label, "render-only bento composition"}
	if len(rows) > p.height {
		rows = rows[:p.height]
	}
	return tea.NewView(strings.Join(rows, "\n"))
}

func viewString(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	if r, ok := v.Content.(interface{ Render() string }); ok {
		return r.Render()
	}
	if s, ok := v.Content.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(v.Content)
}
