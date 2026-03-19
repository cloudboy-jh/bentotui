// Package registry provides access to the embedded brick registry.
package registry

import "embed"

//go:embed bricks/*/* recipes/*/* bentos/*/*
var BricksFS embed.FS
