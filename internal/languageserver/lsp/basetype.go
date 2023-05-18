package lsp

import "net/url"

type (
	Bool        bool
	Decimal     float64
	String      string
	Integer     int32
	UInteger    uint32
	DocumentURI url.URL
	URI         url.URL
)
