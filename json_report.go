package main

import (
    "encoding/json"
    "fmt"
"io/ioutil"
    "os"
    "time"
)

type JSONReport struct {
    Tool        string             `json:"tool"`
    Version     string             `json:"version"`
    TargetURL   string             `json:"target_url"`
    ScanTime    string             `json:"scan_time"`
    Duration    string             `json:"duration"`
    Findings    []JSONFinding      `json:"findings"`
    Exploits    []ExploitResult    `json:"exploits,omitempty"`
}

type JSONFinding struct {
    URL         string `json:"url"`
    Payload     string `json:"payload"`
    VulnType    string `json:"type"`
    Severity    string `json:"severity"`
    Description string `json:"description"`
    Remediation string `json:"remediation"`
}

func ExportJSON(findings []Finding, exploits []ExploitResult, target string, startTime time.Time, filename string) {
    jf := []JSONFinding{}
    for _, f := range findings {
        info := GetVulnerabilityInfo(f.VulnType)
        jf = append(jf, JSONFinding{
            URL:         f.URL,
            Payload:     f.Payload,
            VulnType:    f.VulnType,
            Severity:    f.Severity,
            Description: info.Description,
            Remediation: info.Remediation,
        })
    }
    report := JSONReport{
        Tool:      "ATLAS Scanner V6 ELITE",
        Version:   "6.0",
        TargetURL: target,
        ScanTime:  startTime.Format(time.RFC3339),
        Duration:  time.Since(startTime).String(),
        Findings:  jf,
        Exploits:  exploits,
    }
    data, _ := json.MarshalIndent(report, "", "  ")
    os.MkdirAll("reports", 0755)
    ioutil.WriteFile(filename, data, 0644)
    fmt.Printf("\n[✓] JSON report saved: %s\n", filename)
}
