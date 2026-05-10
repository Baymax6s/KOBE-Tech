package swagger

import "embed"

// Files contains the embedded Swagger UI assets and OpenAPI definition.
//
//go:embed index.html openapi.yml swagger.json
var Files embed.FS
