package iotransport

import (
	"bufio"
	"context"
	"io"
	"net/textproto"

	"github.com/dogmatiq/harpy"
)

type RequestSetReader struct {
	Reader io.Reader
	buf    *bufio.Reader
}

// Read reads the next RequestSet that is to be processed.
//
// It returns ctx.Err() if ctx is canceled while waiting to read the next
// request set. If request set data is read but cannot be parsed a native
// JSON-RPC Error is returned. Any other error indicates an IO error.
func (r *RequestSetReader) Read(ctx context.Context) (harpy.RequestSet, error) {
	if r.buf == nil {
		r.buf = bufio.NewReader(r.Reader)
	}

	reader := textproto.NewReader(r.buf)

	if _, err := reader.ReadMIMEHeader(); err != nil {
		return harpy.RequestSet{}, err
	}

	return harpy.UnmarshalRequestSet(r.buf)
}
