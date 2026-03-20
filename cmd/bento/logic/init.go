package logic

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	bentoregistry "github.com/cloudboy-jh/bentotui/registry"
)

// InstallBento copies a full bento template from the embedded registry.
// The destination directory is created as ./<bento-name>.
func InstallBento(name string) InstallResult {
	result := InstallResult{Name: name}

	if !catalogHasName(BentoRegistry(), name) {
		result.Error = fmt.Errorf("unknown bento: %s", name)
		return result
	}

	destRoot := name
	if _, err := os.Stat(destRoot); err == nil {
		result.Error = fmt.Errorf("directory %q already exists", destRoot)
		return result
	}
	if err := os.MkdirAll(destRoot, 0755); err != nil {
		result.Error = fmt.Errorf("create directory %s: %w", destRoot, err)
		return result
	}

	srcRoot := filepath.ToSlash(filepath.Join("bentos", name))
	walkErr := fs.WalkDir(bentoregistry.BricksFS, srcRoot, func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(srcPath, srcRoot+"/")
		if relPath == "." {
			return nil
		}
		if relPath == srcPath {
			return nil
		}

		dstPath := filepath.Join(destRoot, filepath.FromSlash(relPath))
		if d.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		srcFile, openErr := bentoregistry.BricksFS.Open(srcPath)
		if openErr != nil {
			return openErr
		}

		dstFile, createErr := os.OpenFile(dstPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if createErr != nil {
			srcFile.Close()
			return createErr
		}

		_, copyErr := io.Copy(dstFile, srcFile)
		srcCloseErr := srcFile.Close()
		closeErr := dstFile.Close()
		if copyErr != nil {
			return copyErr
		}
		if srcCloseErr != nil {
			return srcCloseErr
		}
		if closeErr != nil {
			return closeErr
		}

		result.Files = append(result.Files, dstPath)
		return nil
	})

	if walkErr != nil {
		result.Error = fmt.Errorf("install bento %q: %w", name, walkErr)
	}

	return result
}

func catalogHasName(catalog []CatalogEntry, name string) bool {
	for _, c := range catalog {
		if c.Name == name {
			return true
		}
	}
	return false
}
