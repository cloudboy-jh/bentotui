package policy

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestRoomsImportBoundaries(t *testing.T) {
	root := repoRoot(t)
	files := mustGoFiles(t, filepath.Join(root, "registry", "rooms"))

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		imports := mustImports(t, file)
		for _, imp := range imports {
			switch {
			case strings.HasPrefix(imp, "github.com/cloudboy-jh/bentotui/theme"):
				t.Fatalf("rooms must stay theme-agnostic, found %q in %s", imp, rel(root, file))
			case strings.HasPrefix(imp, "github.com/cloudboy-jh/bentotui/registry/bricks/"):
				t.Fatalf("rooms cannot depend on bricks, found %q in %s", imp, rel(root, file))
			case strings.HasPrefix(imp, "charm.land/bubbles"):
				t.Fatalf("rooms cannot depend on bubbles, found %q in %s", imp, rel(root, file))
			}
		}
	}
}

func TestBricksDoNotImportOtherBricks(t *testing.T) {
	root := repoRoot(t)
	files := mustGoFiles(t, filepath.Join(root, "registry", "bricks"))

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		imports := mustImports(t, file)
		for _, imp := range imports {
			if strings.HasPrefix(imp, "github.com/cloudboy-jh/bentotui/registry/bricks/") {
				t.Fatalf("bricks must stay standalone; cross-brick import %q in %s", imp, rel(root, file))
			}
		}
	}
}

func TestBentosAvoidRawBubblesImports(t *testing.T) {
	root := repoRoot(t)
	files := mustGoFiles(t, filepath.Join(root, "registry", "bentos"))
	assertNoRawBubblesImports(t, root, files, "bentos should use Bento bricks")
}

func TestStarterAndScaffoldAvoidRawBubblesImports(t *testing.T) {
	root := repoRoot(t)
	files := append(
		mustGoFiles(t, filepath.Join(root, "cmd", "starter-app")),
		mustGoFiles(t, filepath.Join(root, "cmd", "bento", "logic"))...,
	)
	assertNoRawBubblesImports(t, root, files, "starter/scaffold should keep bubbles behind Bento APIs")
}

func assertNoRawBubblesImports(t *testing.T, root string, files []string, prefix string) {
	t.Helper()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		imports := mustImports(t, file)
		for _, imp := range imports {
			if strings.HasPrefix(imp, "charm.land/bubbles") && imp != "charm.land/bubbles/v2/spinner" {
				t.Fatalf("%s, not raw bubbles import %q in %s", prefix, imp, rel(root, file))
			}
			if strings.HasPrefix(imp, "github.com/charmbracelet/bubbles") && imp != "github.com/charmbracelet/bubbles/spinner" {
				t.Fatalf("%s, not raw bubbles import %q in %s", prefix, imp, rel(root, file))
			}
		}
	}
}

func TestBentoViewsDoNotReadGlobalTheme(t *testing.T) {
	root := repoRoot(t)
	files := mustGoFiles(t, filepath.Join(root, "registry", "bentos"))

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", rel(root, file), err)
		}
		for _, decl := range node.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv == nil || fd.Name == nil || fd.Name.Name != "View" || fd.Body == nil {
				continue
			}
			ast.Inspect(fd.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				x, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}
				if x.Name == "theme" && sel.Sel != nil && sel.Sel.Name == "CurrentTheme" {
					pos := fset.Position(sel.Pos())
					t.Fatalf("View() must use model-owned theme, found theme.CurrentTheme() in %s:%d", rel(root, file), pos.Line)
				}
				return true
			})
		}
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, cur, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve runtime caller path")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(cur), "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("failed to resolve repo root from %s: %v", cur, err)
	}
	return root
}

func mustGoFiles(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".go" {
			out = append(out, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", dir, err)
	}
	return out
}

func mustImports(t *testing.T, file string) []string {
	t.Helper()
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse imports %s: %v", file, err)
	}
	imports := make([]string, 0, len(node.Imports))
	for _, imp := range node.Imports {
		v, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			t.Fatalf("unquote import path %s: %v", imp.Path.Value, err)
		}
		imports = append(imports, v)
	}
	return imports
}

func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(r)
}
