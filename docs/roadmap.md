# BentoTUI Roadmap

## Current State (v0.1)

✅ **Completed:**
- Canvas-based layout system (Horizontal/Vertical with gutters)
- Container components (Panel, Bar, Dialog)
- Widget system with preferred height (Input, List, Text, Card, Table)
- Theme system with 15 presets via bubbletint
- Constraint-based height allocation (Fixed, Min, Max, Flex)

## Next Steps

### v0.2 - Layout Enhancements
- [ ] Grid layout (rows + columns together)
- [ ] Responsive breakpoints (different layouts at different sizes)
- [ ] Layout animations (smooth transitions between states)

### v0.3 - New Widgets
- [ ] Button component
- [ ] Checkbox/Radio components
- [ ] Dropdown/Select component
- [ ] Progress bar
- [ ] Spinner/Loading indicator

### v0.4 - Advanced Features
- [ ] Split pane (draggable resizer between panels)
- [ ] Tabs component
- [ ] Accordion/Collapsible sections
- [ ] Tooltip system

### v0.5 - Developer Experience
- [ ] Component playground/storybook
- [ ] Visual theme editor
- [ ] Layout debugger (show constraints, allocations)
- [ ] Hot reload for development

## Open Questions

1. **Performance:** Canvas-based layout is solid but memory-heavy. Need benchmarks.
2. **Accessibility:** Terminal accessibility (screen readers) needs research.
3. **Mobile:** Does anyone want TUI on mobile? Probably not, but worth considering.

## Immediate Priorities

1. **Bug:** Widget height constraints need validation with complex content
2. **Docs:** Add more usage examples beyond the starter-app
3. **Testing:** Widgets package has no tests

---

Last updated: 2026-03-02
