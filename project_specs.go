package golearning

import _ "embed"

// Project specs live in lessons_mdx so they are easy to edit alongside lessons content.

//go:embed lessons_mdx/Проекты/capstone-rest.md
var CapstoneRESTSpecMD string

//go:embed lessons_mdx/Проекты/capstone-grpc.md
var CapstoneGRPCSpecMD string
