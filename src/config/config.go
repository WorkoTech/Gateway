package config

import (
    "log"
    "net/http"
    "io/ioutil"
    "regexp"
    "strings"

    "worko.tech/gateway/src/utils"

    "gopkg.in/yaml.v2"
)

type AdditionalRequestFields struct {
    WorkspaceId int `json:"workspaceId"`
    FileSize    int `json:"fileSize"`
}

type AccessRuleField struct {
    Name                string          `yaml:"name"`
    RequestValue        string          `yaml:"inRequestValue"`
}

type AccessRule struct {
    Method              string          `yaml:"method"`
    Resource            string          `yaml:"resource"`
    AdditionalFields    []AccessRuleField `yaml:"additionalFields"`
}

type GatewayPath struct {
    Path                string          `yaml:"path"`
    Method              []string        `yaml:"method,flow"`
    Host                string          `yaml:"host"`
    Port                string          `yaml:"port"`
    Protocol            string          `yaml:"protocol"`
    AuthRequired        bool            `yaml:"auth"`
    AccessRules          []AccessRule      `yaml:"accessRule"`
}

type GatewayCfg struct {
    Gateway struct {
        Environment string `yaml:"environment"`
        Port        int `yaml:"port"`
        Paths       []GatewayPath `yaml:"paths"`
    } `yaml:"gateway"`
}

func LoadConfiguration(path string) (*GatewayCfg, error) {
    content, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }

    cfg := GatewayCfg{}
    err = yaml.Unmarshal(content, &cfg)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    return &cfg, nil
}


func (cfg *GatewayCfg) GetPathConfiguration(req *http.Request) *GatewayPath {
    url := req.URL.Path
    if len(url) > 1 && url[len(url) - 1] == '/' {
        url = url[0:len(url) -1]
    }

    for _, path := range cfg.Gateway.Paths {
        match, _ := regexp.MatchString(path.Path, url)
        if match {
            return &path
        }
    }
    return nil
}


func (cfg *GatewayCfg) Log() {
    log.Printf("[INFO] Configuration successfuly loaded")
    log.Printf("[INFO] Environment : %v", utils.GetEnv("ENVIRONMENT", "development"))

    log.Printf("[INFO] Exposed endpoints : ")
    for _, path := range cfg.Gateway.Paths {
        s := "public"
        if path.AuthRequired {
            s = "require auth"
        }
        log.Printf("\t\t - %v %v %v (%v)", path.Method, path.Path, s, strings.ToUpper(path.Protocol))
    }
    log.Printf("")
}

