package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

func FingerprintTech(url string) map[string]string {
    fmt.Println("\n[*] Fingerprinting technologies...")
    client := http.Client{Timeout: 8 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    bodyStr := string(body)
    headers := resp.Header

    techs := make(map[string]string)

    // Server header
    if server := headers.Get("Server"); server != "" {
        techs["Server"] = server
        fmt.Printf("[+] Server: %s\n", server)
    }

    // X-Powered-By
    if powered := headers.Get("X-Powered-By"); powered != "" {
        techs["X-Powered-By"] = powered
        fmt.Printf("[+] Powered By: %s\n", powered)
    }

    // Set-Cookie (PHP session)
    if strings.Contains(headers.Get("Set-Cookie"), "PHPSESSID") {
        techs["PHP"] = "detected"
        fmt.Println("[+] PHP detected (PHPSESSID)")
    }

    // WordPress
    if strings.Contains(bodyStr, "wp-content") || strings.Contains(bodyStr, "wp-includes") {
        techs["WordPress"] = "detected"
        fmt.Println("[+] WordPress detected")
    }

    // jQuery
    if strings.Contains(bodyStr, "jquery") {
        techs["jQuery"] = "detected"
        fmt.Println("[+] jQuery detected")
    }

    // Bootstrap
    if strings.Contains(bodyStr, "bootstrap") {
        techs["Bootstrap"] = "detected"
        fmt.Println("[+] Bootstrap detected")
    }

    // Django
    if strings.Contains(bodyStr, "django") || strings.Contains(bodyStr, "csrftoken") {
        techs["Django"] = "detected"
        fmt.Println("[+] Django detected")
    }

    // Laravel
    if strings.Contains(bodyStr, "laravel") || strings.Contains(headers.Get("Set-Cookie"), "laravel_session") {
        techs["Laravel"] = "detected"
        fmt.Println("[+] Laravel detected")
    }

    return techs
}
