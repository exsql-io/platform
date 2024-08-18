package main

import "embed"

//go:embed third_party/substrait/extensions/*.yaml
var SubstraitExtensions embed.FS
