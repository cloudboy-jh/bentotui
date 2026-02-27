package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id,omitempty"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for s.Scan() {
		line := s.Bytes()
		if len(line) == 0 {
			continue
		}
		var req request
		if err := json.Unmarshal(line, &req); err != nil {
			write(w, response{JSONRPC: "2.0", Error: map[string]any{"code": -32700, "message": err.Error()}})
			continue
		}
		write(w, handle(req))
	}
}

func handle(req request) response {
	switch req.Method {
	case "initialize":
		return response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{
			"serverInfo": map[string]any{"name": "bentotui-mcp", "version": "0.1.0"},
			"capabilities": map[string]any{
				"tools":     map[string]any{},
				"resources": map[string]any{},
			},
		}}
	case "tools/list":
		return response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"tools": []map[string]any{
			{"name": "starter_blueprint", "description": "Return canonical BentoTUI starter app structure"},
			{"name": "component_scaffold", "description": "Return component scaffold contract guidance"},
			{"name": "footer_recipe", "description": "Return footer card contract and truncation recipe"},
			{"name": "header_recipe", "description": "Return header contract mirroring footer behavior"},
			{"name": "focus_recipe", "description": "Return focus API and event contract guidance"},
			{"name": "theme_recipe", "description": "Return theme picker preview/apply/revert flow"},
		}}}
	case "tools/call":
		var p struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(req.Params, &p)
		return response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{
			"content": []map[string]any{{"type": "text", "text": toolText(p.Name)}},
		}}
	case "resources/list":
		return response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"resources": []map[string]any{
			{"uri": "resource://bentotui/contracts", "name": "BentoTUI component contracts"},
			{"uri": "resource://bentotui/starter-app", "name": "Starter app composition rules"},
			{"uri": "resource://bentotui/next-steps", "name": "Immediate execution priorities"},
		}}}
	case "resources/read":
		var p struct {
			URI string `json:"uri"`
		}
		_ = json.Unmarshal(req.Params, &p)
		return response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{
			"contents": []map[string]any{{"uri": p.URI, "text": resourceText(p.URI)}},
		}}
	default:
		return response{JSONRPC: "2.0", ID: req.ID, Error: map[string]any{"code": -32601, "message": fmt.Sprintf("unknown method %s", req.Method)}}
	}
}

func write(w *bufio.Writer, resp response) {
	b, _ := json.Marshal(resp)
	_, _ = w.Write(append(b, '\n'))
	_ = w.Flush()
}

func toolText(name string) string {
	switch name {
	case "starter_blueprint":
		return "Use cmd/starter-app with shell layering header -> body -> footer -> scrim -> dialog. Keep header for context cards and footer for slash command discoverability."
	case "component_scaffold":
		return "Component contract: implement tea.Model plus SetSize/GetSize. Render only within assigned bounds and use semantic styles from ui/styles."
	case "footer_recipe":
		return "Footer contract: one row, segments left|cards|right, truncation priority right>left>cards, card collapse full->command-only->drop from end."
	case "header_recipe":
		return "Header mirrors footer API and behavior, but should carry context/state cards instead of action hints."
	case "focus_recipe":
		return "Focus manager APIs: SetRing, SetIndex, FocusBy, SetEnabled, SetWrap. Emit FocusChangedMsg deterministically when index changes."
	case "theme_recipe":
		return "Theme picker UX: selection movement previews theme (non-persistent), enter commits, esc reverts to pre-open theme and closes."
	default:
		return "Unknown tool"
	}
}

func resourceText(uri string) string {
	switch uri {
	case "resource://bentotui/contracts":
		return "Core contract: tea.Model + SetSize/GetSize, bounded rendering, semantic styles only, deterministic key routing precedence."
	case "resource://bentotui/starter-app":
		return "Starter app is the canonical integration surface for shell, footer/header, dialogs, focus, and theme flow."
	case "resource://bentotui/next-steps":
		return "Immediate priorities are tracked in project-docs/next-steps.md and should be treated as implementation queue."
	default:
		return "Unknown resource"
	}
}
