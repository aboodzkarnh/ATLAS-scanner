package main

import (
    "crypto/rand"
    "math/big"
    "net/http"
    "time"
)

type CustomHeaderConfig struct {
    Headers map[string]string
}

func NewCustomHeaderConfig() *CustomHeaderConfig {
    return &CustomHeaderConfig{
        Headers: make(map[string]string),
    }
}

func (c *CustomHeaderConfig) AddHeader(key, value string) {
    c.Headers[key] = value
}

func (c *CustomHeaderConfig) Apply(req *http.Request) {
    for k, v := range c.Headers {
        req.Header.Set(k, v)
    }
}

type StealthClientWithHeaders struct {
    client     *http.Client
    userAgents []string
    config     Config
    headers    *CustomHeaderConfig
}

func NewStealthClientWithHeaders(cfg Config, headers *CustomHeaderConfig) *StealthClientWithHeaders {
    agents := []string{
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_1) Version/17.1",
    }
    return &StealthClientWithHeaders{
        client:     &http.Client{Timeout: time.Duration(cfg.Stealth.Timeout) * time.Second},
        userAgents: agents,
        config:     cfg,
        headers:    headers,
    }
}

func (s *StealthClientWithHeaders) DoRequest(targetURL string) (*http.Response, error) {
    if s.config.Stealth.DelayMin > 0 {
        n, _ := rand.Int(rand.Reader, big.NewInt(int64(s.config.Stealth.DelayMax-s.config.Stealth.DelayMin)))
        time.Sleep(time.Duration(int(n.Int64())+s.config.Stealth.DelayMin) * time.Millisecond)
    }
    req, _ := http.NewRequest("GET", targetURL, nil)
    if s.config.Stealth.RandomAgent {
        idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s.userAgents))))
        req.Header.Set("User-Agent", s.userAgents[idx.Int64()])
    }
    req.Header.Set("Accept", "*/*")
    if s.headers != nil {
        s.headers.Apply(req)
    }
    return s.client.Do(req)
}
