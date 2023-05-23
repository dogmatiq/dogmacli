package transport

import (
	"context"
	"io"

	"github.com/dogmatiq/harpy"
)

func Run(
	ctx context.Context,
	e harpy.Exchanger,
	r io.Reader,
	w io.Writer,
	l harpy.ExchangeLogger,
) error {
	reader := &RequestSetReader{Reader: r}
	writer := &ResponseWriter{Writer: w}

	for {
		if err := harpy.Exchange(ctx, e, reader, writer, l); err != nil {
			return err
		}
	}
}
