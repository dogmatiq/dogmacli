package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dogmatiq/harpy"
)

// ResponseWriter is an implementation of harpy.ResponseWriter that writes
// responses to an http.ResponseWriter.
type ResponseWriter struct {
	Writer io.Writer

	batch bytes.Buffer
	enc   *json.Encoder
}

// WriteError writes an error response that is a result of some problem with
// the request set as a whole.
//
// It immediately writes the HTTP response headers followed by the HTTP body.
//
// If the error code is pre-defined by the JSON-RPC specification the HTTP
// status code is set to the most appropriate equivalent, otherwise it is set to
// 500 (Internal Server Error).
func (w *ResponseWriter) WriteError(res harpy.ErrorResponse) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return w.flush(data)
}

// WriteUnbatched writes a response to an individual request that was not part
// of a batch.
//
// It immediately writes the HTTP response headers followed by the HTTP body.
//
// If res is an ErrorResponse and its error code is pre-defined by the JSON-RPC
// specification the HTTP status code is set to the most appropriate equivalent.
//
// Application-defined JSON-RPC errors always result in a HTTP 200 (OK), as they
// considered part of normal operation of the transport.
func (w *ResponseWriter) WriteUnbatched(res harpy.Response) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return w.flush(data)
}

// WriteBatched writes a response to an individual request that was part of a
// batch.
//
// If this is the first response of the batch, it immediately writes the HTTP
// response headers and the opening bracket of the array that encapsulates the
// batch of responses.
//
// The HTTP status code is always 200 (OK), as even if res is an ErrorResponse,
// other responses in the batch may indicate a success.
func (w *ResponseWriter) WriteBatched(res harpy.Response) error {
	if w.batch.Len() == 0 {
		w.batch.WriteByte('[')
	} else {
		w.batch.WriteByte(',')
	}

	if w.enc == nil {
		w.enc = json.NewEncoder(&w.batch)
	}

	return w.enc.Encode(res)
}

// Close is called to signal that there are no more responses to be sent.
//
// If batched responses have been written, it writes the closing bracket of the
// array that encapsulates the responses.
func (w *ResponseWriter) Close() error {
	if w.batch.Len() != 0 {
		defer w.batch.Reset()
		w.batch.WriteByte(']')
		return w.flush(w.batch.Bytes())
	}

	return nil
}

func (w *ResponseWriter) flush(data []byte) error {
	if _, err := fmt.Fprintf(
		w.Writer,
		"Content-Length: %d\r\n\r\n",
		len(data),
	); err != nil {
		return err
	}

	_, err := w.Writer.Write(data)
	return err
}
