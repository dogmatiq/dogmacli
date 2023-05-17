package lsp

import "net/url"

type (
	DocumentURI = *url.URL
	URI         = DocumentURI
	Null        struct{}
)
