package main

import (
    "fmt"
    "os"

    "atlas-scanner/scanner"
)

func main() {
    fmt.Println("[ATLAS CMD] Thin CLI wrapper (placeholder)")
    if len(os.Args) < 2 {
        fmt.Println("Usage: atlas <target_url>")
        os.Exit(1)
    }
    target := os.Args[1]
    engine := scanner.NewEngine(target)
    if err := engine.Run(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
