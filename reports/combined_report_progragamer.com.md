# Atlas Scanner Report for https://progragamer.com/

**Scan Date:** 2026-06-20 10:29:53

## Findings Summary

### 1. [High] Admin Panel Bypass
- **URL:** `https://progragamer.com:443/admin/`
- **Payload:** `Accessed /admin/ directly`
- **CVSS:** CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H (8.6 – High)
- **Description:** An administrative login page was discovered. Weak credentials could grant full application control.
- **Remediation:** Restrict access by IP whitelisting. Enforce strong, unique passwords and MFA.

### 2. [High] Admin Panel Bypass
- **URL:** `https://progragamer.com:443/wp-admin/`
- **Payload:** `Accessed /wp-admin/ directly`
- **CVSS:** CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H (8.6 – High)
- **Description:** An administrative login page was discovered. Weak credentials could grant full application control.
- **Remediation:** Restrict access by IP whitelisting. Enforce strong, unique passwords and MFA.

