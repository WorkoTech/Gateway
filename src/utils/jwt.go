package utils

import (
    "net/http"
    "log"
    "strings"

    "github.com/dgrijalva/jwt-go"
)

func extractJwtFromRequest(req *http.Request) (string, bool) {
    // Get "token" from query params
    // Websocket clients (as socket.io) have to send token through query params
    tokens, hasTokenQueryParams := req.URL.Query()["token"]
    if (hasTokenQueryParams) {
        if (len(tokens) > 0 && len(tokens[0]) > 0) {
            return tokens[0], true
        }
        return "", false
    }

    // Get token from headers
    headers := req.Header["Authorization"]
    log.Printf("[DEBUG] headers : %v", headers)
    if len(headers) == 0 {
        return "", false
    }
    return headers[0], true
}

func IsJwtValid(req *http.Request) bool {
    token, err := GetJwtFromRequest(req)

    if err != nil {
        log.Printf("[ERROR] While parsing JWT %v", err)
        return false
    }

    if token == nil {
        return false
    }

    return token.Valid
}

func GetJwtFromRequest(req *http.Request) (*jwt.Token, error) {
    jwtSecret := GetEnv("GATEWAY_JWT_SECRET", "secret")

    // Extract JWT from token
    rawToken, hasToken := extractJwtFromRequest(req)
    if (!hasToken) {
        return nil, nil
    }

    // Remove "Bearer" from rawToken formated as : Bearer <token>
    splittedToken := strings.Split(rawToken, " ")
    if len(splittedToken) != 2 {
        log.Printf("Unable to parse token (%v)", rawToken)
        return nil, nil
    }

    // Parse and verify the token
    return jwt.Parse(splittedToken[1], func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })
}
