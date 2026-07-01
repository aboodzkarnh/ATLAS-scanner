package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/briandowns/spinner"
	"github.com/goccy/go-yaml"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- البنيات الأساسية ---
type Config struct {
	Stealth struct {
		RandomAgent bool `yaml:"random_agent"`
		DelayMin    int  `yaml:"delay_min"`
		DelayMax    int  `yaml:"delay_max"`
		Timeout     int  `yaml:"timeout"`
	} `yaml:"stealth"`
	Scan struct {
		Threads  int `yaml:"threads"`
		MaxDepth int `yaml:"max_depth"`
	} `yaml:"scan"`
	Database struct {
		Type string `yaml:"type"`
		Path string `yaml:"path"`
	} `yaml:"database"`
}

type Scan struct {
	ID        uint      `gorm:"primaryKey"`
	TargetURL string
	StartTime time.Time
	EndTime   time.Time
	Findings  []Finding `gorm:"foreignKey:ScanID"`
}

type Finding struct {
	ID        uint          `gorm:"primaryKey"`
	ScanID    uint
	URL       string
	Payload   string
	VulnType  string
	Severity  string
	Timestamp time.Time
	Exploits  []ExploitResult `gorm:"-"`
}

type Payload struct {
	Name        string   `yaml:"name"`
	Severity    string   `yaml:"severity"`
	Method      string   `yaml:"method"`
	Paths       []string `yaml:"paths"`
	Payloads    []string `yaml:"payloads"`
	DetectRegex string   `yaml:"detect_regex"`
}

type ExploitResult struct {
	Type  string
	Title string
	URL   string
}

type WAFInfo struct {
	Name       string
	BypassTips []string
	IsDetected bool
}

// --- التخفي ---
type StealthClient struct {
	client     *http.Client
	userAgents []string
	config     Config
}

func NewStealthClient(cfg Config) *StealthClient {
	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	}
	return &StealthClient{
		client:     &http.Client{Timeout: time.Duration(cfg.Stealth.Timeout) * time.Second},
		userAgents: agents,
		config:     cfg,
	}
}

func (s *StealthClient) DoRequest(targetURL string) (*http.Response, error) {
	if s.config.Stealth.DelayMin > 0 {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(s.config.Stealth.DelayMax-s.config.Stealth.DelayMin)))
		time.Sleep(time.Duration(int(n.Int64())+s.config.Stealth.DelayMin) * time.Millisecond)
	}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	if s.config.Stealth.RandomAgent {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s.userAgents))))
		req.Header.Set("User-Agent", s.userAgents[idx.Int64()])
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	return s.client.Do(req)
}

// --- الماسح الرئيسي ---
type AtlasScanner struct {
	targetURL   string
	client      *StealthClient
	payloadsDir string
	db          *gorm.DB
	currentScan *Scan
	findings    []Finding
	mu          sync.Mutex
	cfg         Config
}

func NewAtlasScanner(target string, cfg Config) (*AtlasScanner, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Scan{}, &Finding{})
	scan := &Scan{TargetURL: target, StartTime: time.Now()}
	db.Create(scan)
	return &AtlasScanner{
		targetURL:   target,
		client:      NewStealthClient(cfg),
		payloadsDir: "payloads",
		db:          db,
		currentScan: scan,
		cfg:         cfg,
	}, nil
}

func (a *AtlasScanner) LoadPayloads() ([]Payload, error) {
	var all []Payload
	files, err := ioutil.ReadDir(a.payloadsDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".yaml") && !strings.HasSuffix(f.Name(), ".yml") {
			continue
		}
		data, err := ioutil.ReadFile(a.payloadsDir + "/" + f.Name())
		if err != nil {
			continue
		}
		var p Payload
		if err := yaml.Unmarshal(data, &p); err != nil {
			continue
		}
		all = append(all, p)
	}
	return all, nil
}

// --- فلتر WAF الذكي ---
func isBlockedByWAF(body string, statusCode int, header http.Header) bool {
	if header.Get("CF-RAY") != "" || strings.Contains(strings.ToLower(header.Get("Server")), "cloudflare") {
		return true
	}
	bodyLower := strings.ToLower(body)
	if strings.Contains(bodyLower, "attention required") || strings.Contains(bodyLower, "just a moment") {
		return true
	}
	if strings.Contains(bodyLower, "request blocked") || strings.Contains(bodyLower, "access denied") {
		return true
	}
	if statusCode == 403 || statusCode == 503 {
		return true
	}
	return false
}

func isKnownErrorPage(body string) bool {
	lower := strings.ToLower(body)
	errors := []string{
		"page not found", "not found", "error 404", "404 not found",
		"no such file or directory", "bad request", "internal server error",
	}
	for _, e := range errors {
		if strings.Contains(lower, e) {
			return true
		}
	}
	return false
}

// --- محرك الزحف الذكي (Smart Crawler) ---
// يتم استدعاؤه في بداية Run لاستخراج الروابط
func (a *AtlasScanner) crawlSite(maxDepth int) []string {
	fmt.Println("\n[*] Starting smart crawl (extracting internal links)...")
	discovered := make(map[string]bool)
	var toVisit []string
	toVisit = append(toVisit, a.targetURL)
	depth := 0
	maxPages := 50 // حد أقصى للصفحات لتجنب البطء

	for len(toVisit) > 0 && depth < maxDepth {
		if len(discovered) >= maxPages {
			break
		}
		var next []string
		for _, currentURL := range toVisit {
			if discovered[currentURL] {
				continue
			}
			discovered[currentURL] = true
			fmt.Printf("[Crawl] Analyzing: %s\n", currentURL)

			resp, err := a.client.DoRequest(currentURL)
			if err != nil {
				continue
			}
			if resp.StatusCode != 200 {
				resp.Body.Close()
				continue
			}

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}

			doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
				link, exists := s.Attr("href")
				if !exists {
					return
				}
				// تحويل الرابط النسبي إلى مطلق
				base, _ := url.Parse(currentURL)
				full, err := url.Parse(link)
				if err != nil {
					return
				}
				absolute := base.ResolveReference(full)
				absoluteStr := absolute.String()

				// التأكد من أن الرابط ينتمي لنفس النطاق
				if strings.HasPrefix(absoluteStr, a.targetURL) {
					// إزالة التكرار والنصوص غير الضرورية
					if !discovered[absoluteStr] {
						next = append(next, absoluteStr)
					}
				}
			})
		}
		toVisit = next
		depth++
	}

	var result []string
	for u := range discovered {
		result = append(result, u)
	}
	fmt.Printf("[+] Crawling finished. Discovered %d unique pages.\n", len(result))
	return result
}

// --- فحص الثغرات ---
func (a *AtlasScanner) CheckVulnerability(p Payload) {
	for _, path := range p.Paths {
		for _, payloadStr := range p.Payloads {
			testURL := a.targetURL
			if strings.Contains(path, "{PAYLOAD}") {
				testURL = strings.Replace(path, "{PAYLOAD}", url.QueryEscape(payloadStr), -1)
				// التأكد من أن الرابط يبدأ بالبروتوكول
				if !strings.HasPrefix(testURL, "http") {
					testURL = strings.TrimRight(a.targetURL, "/") + "/" + strings.TrimLeft(testURL, "/")
				}
			}

			resp, err := a.client.DoRequest(testURL)
			if err != nil {
				continue
			}
			// فحص سريع لـ WAF
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			bodyStr := string(bodyBytes)

			if isBlockedByWAF(bodyStr, resp.StatusCode, resp.Header) || isKnownErrorPage(bodyStr) {
				continue
			}

			// منطق الاكتشاف
			if p.DetectRegex == "" {
				// فحص Time-based (للـ Blind SQLi)
				start := time.Now()
				resp, _ = a.client.DoRequest(testURL)
				if resp != nil {
					resp.Body.Close()
				}
				elapsed := time.Since(start)
				if elapsed > 4*time.Second {
					finding := Finding{
						ScanID:   a.currentScan.ID,
						URL:      testURL,
						Payload:  payloadStr,
						VulnType: p.Name,
						Severity: p.Severity,
						Timestamp: time.Now(),
					}
					a.mu.Lock()
					a.db.Create(&finding)
					a.findings = append(a.findings, finding)
					a.mu.Unlock()
				}
				continue
			}

			if matched, _ := regexp.MatchString("(?i)"+p.DetectRegex, bodyStr); matched {
				// تحسينات خاصة لـ SSTI
				if p.Name == "Server-Side Template Injection" {
					hostOnly := a.targetURL
					if strings.HasPrefix(hostOnly, "http://") {
						hostOnly = strings.TrimPrefix(hostOnly, "http://")
					}
					if strings.HasPrefix(hostOnly, "https://") {
						hostOnly = strings.TrimPrefix(hostOnly, "https://")
					}
					if idx := strings.Index(hostOnly, "/"); idx != -1 {
						hostOnly = hostOnly[:idx]
					}
					if net.ParseIP(hostOnly) != nil {
						continue // تخطي الـ IPs المباشرة لتجنب الإيجابيات الخاطئة
					}
				}

				finding := Finding{
					ScanID:   a.currentScan.ID,
					URL:      testURL,
					Payload:  payloadStr,
					VulnType: p.Name,
					Severity: p.Severity,
					Timestamp: time.Now(),
				}
				// finding.Exploits = FetchExploits(p.Name) // افترض وجود هذه الدالة أو احذف السطر
				a.mu.Lock()
				a.db.Create(&finding)
				a.findings = append(a.findings, finding)
				a.mu.Unlock()
			}
		}
	}
}

// --- اكتشاف المسارات المخفية ---
var commonPaths = []string{
	"/admin", "/login", "/wp-admin", "/backup", "/api", "/.env",
	"/config", "/dashboard", "/phpinfo.php", "/robots.txt",
	"/sitemap.xml", "/.git/HEAD", "/console", "/swagger",
}

func (a *AtlasScanner) discoverPaths() []string {
	var found []string
	for _, path := range commonPaths {
		fullURL := strings.TrimRight(a.targetURL, "/") + path
		resp, err := a.client.DoRequest(fullURL)
		if err != nil {
			continue
		}
		if resp.StatusCode == 200 {
			found = append(found, path)
		}
		resp.Body.Close()
	}
	return found
}

// --- تحليل WAF ---
func (a *AtlasScanner) detectWAF(resp *http.Response) WAFInfo {
	server := resp.Header.Get("Server")
	cfRay := resp.Header.Get("CF-RAY")
	if strings.Contains(cfRay, "-") || strings.Contains(server, "cloudflare") {
		return WAFInfo{Name: "Cloudflare", BypassTips: []string{"Use WebSocket", "Find Origin IP"}, IsDetected: true}
	}
	if resp.StatusCode == 403 || resp.StatusCode == 406 {
		return WAFInfo{Name: "Generic WAF / IPS", BypassTips: []string{"Double Encoding", "Change GET to POST"}, IsDetected: true}
	}
	return WAFInfo{Name: "Not Detected", IsDetected: false}
}

func (a *AtlasScanner) checkWAF() {
	fmt.Println("\n[*] Analyzing Web Application Firewall (WAF)...")
	resp, err := a.client.DoRequest(a.targetURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	waf := a.detectWAF(resp)
	if waf.IsDetected {
		fmt.Printf("\n[!] WAF DETECTED: %s\n", waf.Name)
		finding := Finding{
			ScanID:    a.currentScan.ID,
			URL:       a.targetURL,
			Payload:   fmt.Sprintf("WAF Detected: %s", waf.Name),
			VulnType:  "WAF/IPS Detection",
			Severity:  "Info",
			Timestamp: time.Now(),
		}
		a.mu.Lock()
		a.db.Create(&finding)
		a.findings = append(a.findings, finding)
		a.mu.Unlock()
	} else {
		fmt.Println("[+] No WAF detected")
	}
}

// --- اكتشاف Origin IP ---
func FindOriginIP(target string) {
	fmt.Println("\n[*] Attempting to find Origin IP...")
	host := target
	if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
	} else if strings.HasPrefix(host, "https://") {
		host = strings.TrimPrefix(host, "https://")
	}
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}
	fmt.Println("[*] Checking crt.sh for direct IPs...")
	crtURL := fmt.Sprintf("https://crt.sh/?q=%s&output=json", host)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(crt
