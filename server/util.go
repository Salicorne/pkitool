package server

import (
	"fmt"
	"net/http"
)

func response(w http.ResponseWriter, statusCode int, err string, args ...interface{}) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf(err, args...)))
}
