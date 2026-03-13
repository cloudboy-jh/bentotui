package layouts

// Pancake:
// +----------------------+
// | header               |
// +----------------------+
// |                      |
// |       content        |
// |                      |
// +----------------------+
// | footer               |
// +----------------------+
// Pancake renders header, content, and footer.
func Pancake(width, height int, header, content, footer Sizable) string {
	return Frame(width, height, header, Static(""), content, footer)
}

// TopbarPancake:
// +----------------------+
// | topbar               |
// +----------------------+
// | header               |
// +----------------------+
// |                      |
// |       content        |
// |                      |
// +----------------------+
// | footer               |
// +----------------------+
// TopbarPancake renders topbar, header, content, and footer.
func TopbarPancake(width, height int, topbar, header, content, footer Sizable) string {
	return Frame(width, height, topbar, header, content, footer)
}
