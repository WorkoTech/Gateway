package utils

import (
    "os"
    "net/http"
    "strings"
)

func Contains(arr []string, str string) bool {
    for _, a := range arr {
        if a == str {
            return true
        }
    }
    return false
}

func GetEnv(key, fallback string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        value = fallback
    }
    return value
}

// IsWebSocketRequest returns a boolean indicating whether the request has the
// headers of a WebSocket handshake request.
func IsWebSocketRequest(r *http.Request) bool {
    headerContains := func(key, val string) bool {
        vv := strings.Split(r.Header.Get(key), ",")
        for _, v := range vv {
            if val == strings.ToLower(strings.TrimSpace(v)) {
                return true
            }
        }
        return false
    }

    if !headerContains("Connection", "upgrade") {
        return false
    }

    if !headerContains("Upgrade", "websocket") {
        return false
    }
    return true
}
