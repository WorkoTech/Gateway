package main

import (
    "os"
    "log"
    "net/http"

    "worko.tech/gateway/src/config"
    "worko.tech/gateway/src/handlers"
    "worko.tech/gateway/src/utils"
)

func main() {
    logStart()

    // Get current dir
    dir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    // Load and log configuration
    cfg, err := config.LoadConfiguration(dir + "/cfg_gateway.yaml")
    if err != nil {
        log.Fatal("[ERROR] Unable to load configuration, exiting.")
        return
    }
    cfg.Log()

    // Configure web server
    mux := http.NewServeMux()
    mux.HandleFunc("/health", handlers.Health)
    mux.HandleFunc("/", handlers.ReverseProxy(cfg))

    // Start web server
    port := "0.0.0.0:" + utils.GetEnv("GATEWAY_PORT", "80")
    log.Printf("Gateway starting listening on [:" + port + "] ...\n")
    log.Fatal(http.ListenAndServe(port, mux))
}

func logStart() {
    log.Printf("      ___           ___           ___           ___           ___     ")
    log.Printf("     /__/\\         /  /\\         /  /\\         /__/|         /  /\\    ")
    log.Printf("    _\\_ \\:\\       /  /::\\       /  /::\\       |  |:|        /  /::\\   ")
    log.Printf("   /__/\\ \\:\\     /  /:/\\:\\     /  /:/\\:\\      |  |:|       /  /:/\\:\\  ")
    log.Printf("  _\\_ \\:\\ \\:\\   /  /:/  \\:\\   /  /:/~/:/    __|  |:|      /  /:/  \\:\\ ")
    log.Printf(" /__/\\ \\:\\ \\:\\ /__/:/ \\__\\:\\ /__/:/ /:/___ /__/\\_|:|____ /__/:/ \\__\\:\\")
    log.Printf(" \\  \\:\\ \\:\\/:/ \\  \\:\\ /  /:/ \\  \\:\\/:::::/ \\  \\:\\/:::::/ \\  \\:\\ /  /:/")
    log.Printf("  \\  \\:\\ \\::/   \\  \\:\\  /:/   \\  \\::/~~~~   \\  \\::/~~~~   \\  \\:\\  /:/ ")
    log.Printf("   \\  \\:\\/:/     \\  \\:\\/:/     \\  \\:\\        \\  \\:\\        \\  \\:\\/:/  ")
    log.Printf("    \\  \\::/       \\  \\::/       \\  \\:\\        \\  \\:\\        \\  \\::/   ")
    log.Printf("     \\__\\/         \\__\\/         \\__\\/         \\__\\/         \\__\\/    ")
    log.Printf("")
    log.Printf("[INFO] Starting Worko Gateway...")
}
