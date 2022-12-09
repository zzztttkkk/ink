package ink

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int
	header     http.Header
	buf        *bytes.Buffer
}

func (resp *Response) Header() http.Header {
	if resp.header == nil {
		resp.header = map[string][]string{}
	}
	return resp.header
}

func (resp *Response) Write(p []byte) (int, error) {
	if resp.buf == nil {
		resp.buf = responseWBufPool.Get().(*bytes.Buffer)
	}
	return resp.buf.Write(p)
}

func (resp *Response) WriteString(v string) (int, error) {
	if resp.buf == nil {
		resp.buf = responseWBufPool.Get().(*bytes.Buffer)
	}
	return resp.buf.WriteString(v)
}

func (resp *Response) WriteJSON(v any) error {
	if resp.buf == nil {
		resp.buf = responseWBufPool.Get().(*bytes.Buffer)
	}
	resp.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(resp.buf)
	return encoder.Encode(v)
}

func (resp *Response) Len() int { return resp.buf.Len() }

func (resp *Response) Cap() int { return resp.buf.Cap() }

func (resp *Response) Grow(v int) { resp.buf.Grow(v) }

func (resp *Response) Truncate(v int) { resp.buf.Truncate(v) }

func (resp *Response) Reset() { resp.buf.Reset() }

func (resp *Response) ResetAll() {
	resp.StatusCode = 0
	resp.header = nil
	if resp.buf != nil {
		resp.buf.Reset()
	}
}
