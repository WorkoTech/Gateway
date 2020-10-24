package handlers

import (
    "bytes"
    "log"
    "reflect"
    "encoding/json"
    "io/ioutil"
    "net/url"
    "net/http"
    "net/http/httputil"

    "worko.tech/gateway/src/config"
    "worko.tech/gateway/src/externals"
    "worko.tech/gateway/src/utils"
    "worko.tech/gateway/src/wsutil"
)

func ReverseProxy(cfg *config.GatewayCfg) (func(http.ResponseWriter, *http.Request)) {
    return func(w http.ResponseWriter, req *http.Request) {
        path := cfg.GetPathConfiguration(req)

        if req.Method == "OPTIONS" || req.Method == "HEAD" {
            utils.SetCommonHeaders(w, http.StatusOK)
            return
        }

        if path == nil {
            log.Printf("[WARNING] Not found on %v", req.URL.Path)
            utils.SetCommonHeaders(w, http.StatusNotFound)
            return
        }

        if utils.IsWebSocketRequest(req) {
            websocketHandler(w, path, req)
            return
        }
        httpHandler(w, path, req)
        return
    }
}

func httpHandler(w http.ResponseWriter, path *config.GatewayPath, req *http.Request) {
    // Check if method is allowed in cfg
    if !utils.Contains(path.Method, req.Method) {
        utils.SetCommonHeaders(w, http.StatusMethodNotAllowed)
        log.Printf("[WARNING] Wrong Method use on %v (%v)", path.Path, req.Method)
        return
    }
    // Check auth
    if path.AuthRequired {
        if !utils.IsJwtValid(req) {
            utils.SetCommonHeaders(w, http.StatusUnauthorized)
            log.Printf("[WARNING] Unauthorized connexion on (%v) %v (http)", req.Method, path.Path)
            return
        }

        accessGranted, err := assertAccess(path.AccessRules, req)
        if err != nil {
            log.Printf("[ERROR] Error while checking access : %v", err.Error())
        }
        if !accessGranted {
            log.Printf("[WARNING] Acces denied")
            utils.SetCommonHeaders(w, http.StatusForbidden)
            return
        }
    }

    // Build target host url
    url, _ := url.Parse("http://" + path.Host + ":" + path.Port + "/")
    log.Printf("Serving %v", url)

    // Send the request to target
    proxy := httputil.NewSingleHostReverseProxy(url)
    proxy.ServeHTTP(w, req)
}

func websocketHandler(w http.ResponseWriter, path *config.GatewayPath, req *http.Request) {
    if path.Protocol != "websocket" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // Check Auth
    if path.AuthRequired && !utils.IsJwtValid(req) {
        w.WriteHeader(http.StatusUnauthorized)
        log.Printf("[WARNING] Unauthorized connexion on %v (websocket)", path.Path)
        return
    }

    // Build target host url
    url, _ := url.Parse("ws://" + path.Host + ":" + path.Port)
    log.Printf("Proxying WS req to %v", url)

    // ServeHTTP add query params from req to url
    proxy := wsutil.NewSingleHostWsReverseProxy(url)
    proxy.ServeHTTP(w, req)
}

func assertAccess(accessRules []config.AccessRule, req *http.Request) (bool, error) {
    // If no access rule for this path, grant access
    var accessRule config.AccessRule
    for _, rule := range accessRules {
        if rule.Method == req.Method {
            accessRule = rule
            break
        }
    }
    log.Printf("Access rule : %v", accessRule)
    if accessRule.Resource == "" {
        return true, nil
    }

    // Extract JWT from request
    token, err := utils.GetJwtFromRequest(req)
    if err != nil {
        return false, nil
    }

    // Reading request body and retrieve additionalFields
    // defined in config from request to send them to iam
    // for grant control
    bodyBytes, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Printf("error : %v", err)
        return false, err
    }
    var requestBodyPayload config.AdditionalRequestFields
    err = json.Unmarshal(bodyBytes, &requestBodyPayload)

    var additionalFields = map[string]interface{} {}
    for _, field := range accessRule.AdditionalFields {
        r := reflect.ValueOf(requestBodyPayload)
        f := reflect.Indirect(r).FieldByName(field.RequestValue)

        additionalFields[field.Name] = int(f.Int())
    }

    req.Body.Close()
    req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

    // Get action from request method
    action := ""
    switch req.Method {
        case "POST":
            action = "create"
            break
        case "GET":
            action = "read"
            break
        case "PUT":
            action = "update"
            break
        case "DELETE":
            action = "delete"
            break
    }

    // Post to IAM
    granted, err := externals.RetrieveAccess(token, action, accessRule.Resource, additionalFields)
    if err != nil {
        return false, nil
    }

    return granted, nil
}
