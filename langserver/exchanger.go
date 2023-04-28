package langserver

import (
	"context"

	"github.com/dogmatiq/harpy"
)

type exchanger struct{}

// Call handles call request and returns its response.
func (e *exchanger) Call(ctx context.Context, req harpy.Request) harpy.Response {
	return harpy.NewErrorResponse(
		req.ID,
		harpy.MethodNotFound(),
	)
}

// Notify handles a notification request, which does not expect a response.
func (e *exchanger) Notify(context.Context, harpy.Request) {
}
