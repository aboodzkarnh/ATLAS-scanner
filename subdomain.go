package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

func DiscoverSubdomains(domain string) []string {
    fmt.Println("\n[*] Discovering subdomains via crt.sh...")
    url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
    client := http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    var certs []struct {
        NameValue string `json:"name_value"`
    }
    if err := json.Unmarshal(body, &certs); err != nil {
        return nil
    }
    seen := make(map[string]bool)
    var subs []string
    for _, c := range certs {
        names := strings.Split(c.NameValue, "\n")
        for _, name := range names {
            name = strings.TrimSpace(name)
            if name != "" && !seen[name] && strings.HasSuffix(name, "."+domain) {
                seen[name] = true
                subs = append(subs, name)
                fmt.Printf("[+] Found subdomain: %s\n", name)
            }
        }
    }
    return subs
}
