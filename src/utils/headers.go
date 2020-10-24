package utils

import (
    "net/http"
)

func SetCommonHeaders(w http.ResponseWriter, statusCode int) {
    w.Header().Set("Access-Control-Allow-Methods", GetEnv("ALLOW_METHODS", "GET, POST, PUT, DELETE, OPTIONS, HEAD"))
    w.Header().Set("Access-Control-Allow-Origin", GetEnv("ALLOW_ORIGIN", "*"))
    w.Header().Set("Access-Control-Allow-Headers", GetEnv("ALLOW_HEADERS", "*"))
    w.WriteHeader(statusCode)
}
