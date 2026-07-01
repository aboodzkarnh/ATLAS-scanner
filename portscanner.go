package main

import (
    "fmt"
    "net"
    "strconv"
    "time"
)

var CommonPorts = []int{80, 443, 8080, 8443, 8000, 8888, 3000, 5000, 9090, 9091, 9443, 10443}

func ScanPorts(host string, ports []int) []int {
    var openPorts []int
    fmt.Println("\n[*] Starting stealth port scan...")
    for _, port := range ports {
        address := net.JoinHostPort(host, strconv.Itoa(port))
        conn, err := net.DialTimeout("tcp", address, 2*time.Second)
        if err == nil {
            openPorts = append(openPorts, port)
            fmt.Printf("[+] Port %d is OPEN\n", port)
            conn.Close()
        }
    }
    return openPorts
}
