package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

type TrialCounter struct { filename string }
func NewTrialCounter() *TrialCounter { return &TrialCounter{filename: ".atlas_trial"} }
func (t *TrialCounter) GetCount() int {
    data, err := ioutil.ReadFile(t.filename)
    if err != nil { ioutil.WriteFile(t.filename, []byte("3"), 0644); return 3 }
    count, _ := strconv.Atoi(strings.TrimSpace(string(data)))
    return count
}
func (t *TrialCounter) Decrement() int {
    count := t.GetCount()
    if count > 0 { count--; ioutil.WriteFile(t.filename, []byte(strconv.Itoa(count)), 0644) }
    return count
}
func getMachineID() string {
    cmd := exec.Command("sh", "-c", "dmidecode -s system-uuid 2>/dev/null || echo ''")
    output, err := cmd.Output()
    if err == nil && strings.TrimSpace(string(output)) != "" { return strings.TrimSpace(string(output)) }
    data, err := ioutil.ReadFile("/etc/machine-id")
    if err == nil && strings.TrimSpace(string(data)) != "" { return strings.TrimSpace(string(data)) }
    hostname, _ := os.Hostname(); return hostname
}
func generateLicenseSignature(subscriberID, machineID, expiry string) string {
    secretKey := "ATLAS_SECRET_KEY_2026"
    data := subscriberID + machineID + expiry
    mac := hmac.New(sha256.New, []byte(secretKey))
    mac.Write([]byte(data))
    return hex.EncodeToString(mac.Sum(nil))[:8]
}
func ValidateLicense(key string) (bool, time.Time, string) {
    parts := strings.Split(key, "-")
    if len(parts) != 4 { return false, time.Time{}, "Invalid format" }
    subscriberID, expiryStr, machinePartial, signature := parts[0], parts[1], parts[2], parts[3]
    expiry, err := time.Parse("20060102", expiryStr)
    if err != nil { return false, time.Time{}, "Invalid date" }
    if time.Now().After(expiry) { return false, expiry, "Expired" }
    currentMachineID := getMachineID()
    expectedSig := generateLicenseSignature(subscriberID, currentMachineID, expiryStr)
    if signature == expectedSig { return true, expiry, subscriberID }
    expectedSig = generateLicenseSignature(subscriberID, machinePartial, expiryStr)
    if signature == expectedSig { return true, expiry, subscriberID }
    return false, expiry, "Invalid signature"
}
func CheckLicenseFile() (bool, bool) {
    licenseFile := "license.key"
    trial := NewTrialCounter()
    data, err := ioutil.ReadFile(licenseFile)
    if err == nil {
        licenseKey := strings.TrimSpace(string(data))
        valid, expiry, subscriberID := ValidateLicense(licenseKey)
        if valid {
            daysLeft := int(time.Until(expiry).Hours() / 24)
            fmt.Printf("\n[✓] Licensed to: %s\n", subscriberID)
            fmt.Printf("[✓] Expires: %s (%d days remaining)\n", expiry.Format("January 2, 2006"), daysLeft)
            if daysLeft <= 7 { fmt.Println("[!] Your license expires soon. Please renew.") }
            return true, true
        }
        fmt.Printf("\n[!] License invalid or expired.\n")
    }
    remaining := trial.GetCount()
    if remaining > 0 {
        fmt.Printf("\n[✓] TRIAL MODE - %d runs remaining\n", remaining)
        trial.Decrement()
        return true, false
    }
    fmt.Println("\n[!] Trial expired. No runs remaining.")
    fmt.Println("[*] Your Machine ID:", getMachineID())
    return false, false
}
