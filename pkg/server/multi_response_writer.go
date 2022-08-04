package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type ResponseWriterRecorder struct {
	recorder *httptest.ResponseRecorder
	original http.ResponseWriter
}

func (mw ResponseWriterRecorder) Header() http.Header {
	return mw.original.Header()
}

func (mw ResponseWriterRecorder) Write(bytes []byte) (int, error) {
	mw.recorder.Write(bytes)
	return mw.original.Write(bytes)
}

func (mw ResponseWriterRecorder) WriteHeader(statusCode int) {
	mw.recorder.WriteHeader(statusCode)
	mw.original.WriteHeader(statusCode)
}

func (mw ResponseWriterRecorder) Body() string {
	b, _ := ioutil.ReadAll(mw.recorder.Body)
	return string(b)
}

func NewResponseWriterRecorder(original http.ResponseWriter) *ResponseWriterRecorder {
	return &ResponseWriterRecorder{httptest.NewRecorder(), original}
}
