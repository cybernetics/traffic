package traffic

import (
  "fmt"
  "net/http"
  "encoding/json"
  "encoding/xml"
)

type ResponseWriter interface {
  http.ResponseWriter
  SetVar(string, interface{})
  GetVar(string) interface{}
  StatusCode() int
  Written() bool
  RenderTemplate(string, ...interface{})
  RenderJSON(data interface{})
  RenderXML(data interface{})
  RenderText(string, ...interface{})
}

type responseWriter struct {
  http.ResponseWriter
  written             bool
  bodyWritten         bool
  statusCode          int
  env                 map[string]interface{}
  routerEnv           *map[string]interface{}
  beforeWriteHandlers []func()
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
  w.written     = true
  w.bodyWritten = true

  return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteHeader(statusCode int) {
  w.statusCode = statusCode
  w.ResponseWriter.WriteHeader(statusCode)
  w.written = true
}

func (w *responseWriter) StatusCode() int {
  return w.statusCode
}

func (w *responseWriter) SetVar(key string, value interface{}) {
  w.env[key] = value
}

func (w *responseWriter) Written() bool {
  return w.written
}

func (w *responseWriter) GetVar(key string) interface{} {
  // local env
  value := w.env[key]
  if value != nil {
    return value
  }

  // router env
  value = (*w.routerEnv)[key]
  if value != nil {
    return value
  }

  // global env
  return GetVar(key)
}

func (w *responseWriter) RenderTemplate(templateName string, data ...interface{}) {
  RenderTemplate(w, templateName, data...)
}

func (w *responseWriter) RenderJSON(data interface{}) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  json.NewEncoder(w).Encode(data)
}

func (w *responseWriter) RenderXML(data interface{}) {
  w.Header().Set("Content-Type", "application/xml; charset=utf-8")
  xml.NewEncoder(w).Encode(data)
}

func (w *responseWriter) RenderText(textFormat string, data ...interface{}) {
  fmt.Fprintf(w, textFormat, data...)
}

func newResponseWriter(w http.ResponseWriter, routerEnv *map[string]interface{}) *responseWriter {
  rw := &responseWriter{
    ResponseWriter:       w,
    statusCode:           http.StatusOK,
    env:                  make(map[string]interface{}),
    routerEnv:            routerEnv,
    beforeWriteHandlers:  make([]func(), 0),
  }

  return rw
}

