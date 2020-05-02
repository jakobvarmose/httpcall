package httpcall

import (
	"io"
	"net/http"
)

// Call calls the handler's ServeHTTP method with the given request
// and returns the response.
func Call(wrap http.Handler, req *http.Request) *http.Response {
	bodyReader, bodyWriter := io.Pipe()
	res := &http.Response{
		Header:        make(http.Header),
		Body:          bodyReader,
		ContentLength: -1,
	}
	ch := make(chan struct{})
	w := &writer{
		res:  res,
		ch:   ch,
		body: bodyWriter,
	}
	go func() {
		defer func() {
			if w.ch != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			bodyWriter.Close()
		}()
		wrap.ServeHTTP(w, req)
	}()
	go func() {
		<-req.Context().Done()
		bodyReader.Close()
	}()
	<-ch
	return res
}

// Write writes the response to the response writer.
func Write(w http.ResponseWriter, res *http.Response) error {
	header := w.Header()
	for name, values := range res.Header {
		header[name] = values
	}
	w.WriteHeader(res.StatusCode)
	_, err := io.Copy(w, res.Body)
	return err
}

type writer struct {
	res  *http.Response
	ch   chan<- struct{}
	body *io.PipeWriter
}

func (w *writer) Write(b []byte) (int, error) {
	if w.ch != nil {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(b)
}
func (w *writer) WriteHeader(statusCode int) {
	w.res.StatusCode = statusCode
	close(w.ch)
	w.ch = nil
}
func (w *writer) Header() http.Header {
	return w.res.Header
}
