package handlers

import (
    "io"
    "net/http"
)


func Health(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "OK")
}
