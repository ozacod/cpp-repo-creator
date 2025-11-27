package embedded

import "embed"

//go:embed recipes/*.yaml
var RecipesFS embed.FS
