package main

type VulnerabilityInfo struct {
    Name        string
    Description string
    Severity    string
    Remediation string
}

var knowledgeBase = map[string]VulnerabilityInfo{
    "SQL Injection": {
        Name: "SQL Injection", Severity: "Critical",
        Description: "Attackers inject malicious SQL code to manipulate the database, potentially stealing, modifying, or deleting data.",
        Remediation: "Use parameterized queries (prepared statements) for all database access. Validate and sanitize all user inputs.",
    },
    "Cross-Site Scripting (XSS)": {
        Name: "Cross-Site Scripting (XSS)", Severity: "Medium",
        Description: "Attackers inject client-side scripts into pages viewed by others, leading to session hijacking or defacement.",
        Remediation: "Implement proper output encoding (HTML escaping). Use Content-Security-Policy (CSP) headers.",
    },
    "Local File Inclusion": {
        Name: "Local File Inclusion", Severity: "High",
        Description: "Allows reading of local server files (e.g., /etc/passwd), potentially exposing sensitive data or enabling code execution.",
        Remediation: "Avoid including files based on user input. Use a whitelist of allowed filenames.",
    },
    "Command Injection": {
        Name: "Command Injection", Severity: "Critical",
        Description: "Attackers execute arbitrary system commands on the server, potentially gaining full control.",
        Remediation: "Never pass user input to system commands. Use library functions instead.",
    },
    "Open Redirect": {
        Name: "Open Redirect", Severity: "Low",
        Description: "Web application redirects users to an attacker-specified external URL, often used in phishing.",
        Remediation: "Avoid user-supplied redirect URLs. Use a whitelist of allowed destinations.",
    },
    "Admin Panel Bypass": {
        Name: "Admin Panel Bypass", Severity: "High",
        Description: "An administrative login page was discovered. Weak credentials could grant full application control.",
        Remediation: "Restrict access by IP whitelisting. Enforce strong, unique passwords and MFA.",
    },
    "WAF/IPS Detection": {
        Name: "WAF/IPS Detection", Severity: "Info",
        Description: "A Web Application Firewall was identified. While protective, misconfiguration may allow bypass.",
        Remediation: "Regularly update WAF rules and monitor logs for bypass attempts.",
    },
    "Server-Side Template Injection": {
        Name: "Server-Side Template Injection", Severity: "High",
        Description: "Injects template code that executes server-side, potentially leading to remote code execution.",
        Remediation: "Sanitize user input before template processing. Use sandboxed template engines.",
    },
    "Directory Traversal": {
        Name: "Directory Traversal", Severity: "Medium",
        Description: "Accesses files and directories outside the web root, exposing sensitive system files.",
        Remediation: "Validate and sanitize file paths. Use a whitelist of allowed directories.",
    },
    "XML External Entity": {
        Name: "XML External Entity", Severity: "High",
        Description: "Exploits vulnerable XML parsers to read local files, perform SSRF, or cause denial of service.",
        Remediation: "Disable external entity processing in XML parsers. Use JSON where possible.",
    },
    "Insecure Direct Object Reference": {
        Name: "Insecure Direct Object Reference", Severity: "Medium",
        Description: "Access controls are bypassed by modifying object references (e.g., user IDs in URLs).",
        Remediation: "Implement proper access control checks for every request. Use indirect object references.",
    },
    "Cross-Site Request Forgery": {
        Name: "Cross-Site Request Forgery", Severity: "Medium",
        Description: "Forces authenticated users to perform unwanted actions without their knowledge.",
        Remediation: "Implement anti-CSRF tokens in all state-changing requests.",
    },
    "Exposed .git Directory": {
        Name: "Exposed .git Directory", Severity: "High",
        Description: "Git repository metadata is publicly accessible, potentially leaking source code and secrets.",
        Remediation: "Restrict access to .git directories at the web server level.",
    },
    "PHP Info Disclosure": {
        Name: "PHP Info Disclosure", Severity: "Low",
        Description: "PHP configuration page (phpinfo) is exposed, revealing server paths, versions, and settings.",
        Remediation: "Remove or restrict access to phpinfo files in production environments.",
    },
    "Sensitive Path in robots.txt": {
        Name: "Sensitive Path in robots.txt", Severity: "Low",
        Description: "The robots.txt file reveals hidden paths (e.g., /admin) intended to be kept secret.",
        Remediation: "Do not rely on robots.txt for security. Implement proper access controls instead.",
    },
    "Environment File Leak": {
        Name: "Environment File Leak", Severity: "High",
        Description: "Environment configuration files (.env) are exposed, leaking database credentials and API keys.",
        Remediation: "Block access to dotfiles at the web server level.",
    },
    "Swagger UI Exposure": {
        Name: "Swagger UI Exposure", Severity: "Medium",
        Description: "API documentation (Swagger UI) is publicly accessible, revealing API endpoints and parameters.",
        Remediation: "Restrict access to API documentation to authorized users only.",
    },
}

func GetVulnerabilityInfo(vulnType string) VulnerabilityInfo {
    if info, exists := knowledgeBase[vulnType]; exists { return info }
    return VulnerabilityInfo{
        Name: vulnType, Severity: "Unknown",
        Description: "Detailed information not available for this vulnerability.",
        Remediation: "Consult security best practices for remediation guidance.",
    }
}

// getCVSS returns the CVSS v3.1 vector string and score for a vulnerability type
func getCVSS(vulnType string) string {
    switch vulnType {
    case "SQL Injection", "Command Injection":
        return "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:C/C:H/I:H/A:H (9.8 – Critical)"
    case "Server-Side Template Injection":
        return "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:C/C:H/I:H/A:H (9.8 – Critical)"
    case "Local File Inclusion", "Admin Panel Bypass", "XML External Entity", "Exposed .git Directory":
        return "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H (8.6 – High)"
    case "Cross-Site Scripting (XSS)":
        return "CVSS:3.1/AV:N/AC:L/PR:N/UI:R/S:C/C:L/I:L/A:N (6.1 – Medium)"
    case "Open Redirect":
        return "CVSS:3.1/AV:N/AC:L/PR:N/UI:R/S:C/C:L/I:L/A:N (5.4 – Low)"
    default:
        return "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:N (Score varies)"
    }
}
