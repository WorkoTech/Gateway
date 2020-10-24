package externals

import (
    "bytes"
    "os"
    "log"
    "net/http"
    "encoding/json"

    "github.com/dgrijalva/jwt-go"
)

func RetrieveAccess(token *jwt.Token, action string, resource string, additionalFields map[string]interface{}) (bool, error) {
    log.Printf("RetrieveAccess %v | %v | %v", token, resource, additionalFields)

    client := &http.Client{}

    payload := map[string]interface{}{"resource": resource, "action": action}
    for k, v := range additionalFields {
        payload[k] = v
    }
    jsonPayload, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST","http://" + os.Getenv("IAM_HOST") + ":" + os.Getenv("IAM_PORT") + "/access",bytes.NewBuffer(jsonPayload))
    if err != nil {
        log.Printf("Error : %v", err.Error())
        return false, err
    }

    log.Printf("token %v", token)
    req.Header.Add("Authorization", "Bearer " + token.Raw)
    req.Header.Add("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error : %v", err.Error())
        return false, err
    }
    defer resp.Body.Close()

    return resp.StatusCode == http.StatusOK, nil
}
