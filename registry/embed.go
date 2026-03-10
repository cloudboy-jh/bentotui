// Package registry provides access to the embedded component registry.
package registry

import "embed"

//go:embed components/*/* bentos/*/*
var ComponentsFS embed.FS
