package templates

import (
	"embed"
	_ "embed"
)

//go:embed generator/method_interception.gotmpl
var ProxyTemplate string

//go:embed generator/mocks.gotmpl
var MockTemplate string

//go:embed bootstrap/*
var BootstrapTemplates embed.FS
