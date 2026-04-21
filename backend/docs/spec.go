package docs

import _ "embed"

// SwaggerYAML is the OpenAPI 3 spec served at runtime (embedded so containers and
// arbitrary working directories do not need a docs/ folder on disk).
//
//go:embed swagger.yaml
var SwaggerYAML []byte
