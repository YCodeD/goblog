package middlewares

import "net/http"

// HttpHandlerFunc 简写 --func(w http.ResponseWriter, r *http.Request)
type HttpHandlerFunc func(w http.ResponseWriter, r *http.Request)
