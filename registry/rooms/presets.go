package rooms

// AppShell:
// +-------------------------+
// |         content         |
// +-------------------------+
// |        footer           |
// +-------------------------+
// AppShell renders the common terminal app shape used by command tools.
func AppShell(width, height int, content, footer Sizable) string {
	return Focus(width, height, content, footer)
}

// SidebarDetail:
// +--------+----------------+
// | sidebar|     detail     |
// |        |                |
// +-------------------------+
// |        footer           |
// +-------------------------+
// SidebarDetail renders a two-pane details layout with a footer lane.
func SidebarDetail(width, height, sidebarWidth int, sidebar, detail, footer Sizable) string {
	body := Rail(width, max(1, height-1), sidebarWidth, sidebar, detail)
	return Focus(width, height, Static(body), footer)
}

// Dashboard:
// +-----------+-------------+
// |     tl    |     tr      |
// +-----------+-------------+
// |     bl    |     br      |
// +-------------------------+
// |         footer          |
// +-------------------------+
// Dashboard renders the standard 2x2 metric layout plus footer.
func Dashboard(width, height int, tl, tr, bl, br, footer Sizable) string {
	return Dashboard2x2Footer(width, height, tl, tr, bl, br, footer)
}

// DiffWorkspace:
// +-------------------------+
// |         header          |
// +--------+----------------+
// | files  |      diff      |
// | rail   |      pane      |
// +-------------------------+
// |         footer          |
// +-------------------------+
// DiffWorkspace renders a high-level diff viewer room contract.
func DiffWorkspace(width, height, railWidth int, header, fileRail, diffMain, footer Sizable) string {
	return HolyGrail(width, height, railWidth, header, fileRail, diffMain, footer)
}
