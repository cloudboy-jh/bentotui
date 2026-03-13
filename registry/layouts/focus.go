package layouts

// Focus:
// +----------------------+
// |                      |
// |       content        |
// |                      |
// +----------------------+
// | footer               |
// +----------------------+
// Focus renders full content plus a one-row footer.
func Focus(width, height int, content, footer Sizable) string {
	return Frame(width, height, Static(""), Static(""), content, footer)
}
