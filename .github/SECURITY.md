# Security Policy

The IDXStock team takes the security of our application, dependencies, and data feeds seriously. We appreciate the responsible disclosure of any vulnerabilities found within this project.

---

## Supported Versions

Security fixes and dependency patches are actively maintained for the latest code on the default branch and current minor releases:

| Version / Branch | Supported          | Security Notes |
| :--- | :---: | :--- |
| `main` (active) | :white_check_mark: | Receives immediate security updates and dependency bumps |
| `>= v1.0.0` | :white_check_mark: | Supported for critical CVE patches |
| `< v1.0.0` | :x: | End of Life (EOL) — please upgrade to the latest release |

---

## Reporting a Vulnerability

If you discover a security vulnerability or potential threat in IDXStock, **please do not disclose it publicly** via GitHub Issues, Discussions, or pull requests.

Instead, report it privately through one of the following channels:

### 1. GitHub Private Vulnerability Reporting (Recommended)
Use GitHub's built-in private advisory mechanism:
- Navigate to the **[Security Advisory](https://github.com/KAnggara75/IDXStock/security/advisories)** tab of this repository.
- Click **"Report a vulnerability"** to open a confidential report directly with repository maintainers.

### 2. Direct Security Contact
If you cannot use GitHub Security Advisories, send an encrypted or direct email to:
- **Email**: [kanggara75@gmail.com](mailto:kanggara75@gmail.com)
- **Subject**: `[SECURITY VULNERABILITY] IDXStock - <Brief Title>`

---

## What to Include in Your Report

To help us investigate, triage, and patch the issue promptly, please include:
- **Component / Endpoint**: Specific file, route, handler, or service affected (e.g., `/api/v1/stocks/upload`, SQL queries, or outbound HTTP requests).
- **Vulnerability Type**: E.g., SQL Injection, Remote Code Execution (RCE), Authentication Bypass, Denial of Service (DoS), Insecure Deserialization, SSRF.
- **Severity & Impact**: Your assessment of the practical severity and threat model.
- **Proof of Concept (PoC)**: Detailed step-by-step instructions, HTTP payload samples, or script reproducing the behavior.
- **Mitigation / Suggested Fix**: Any proposed patches or mitigations (optional).

---

## Response & Disclosure Process

1. **Initial Acknowledgment**: We will acknowledge receipt of your vulnerability report within **48 hours**.
2. **Investigation & Triage**: Maintainers will review the reproduction steps, assess the severity, and confirm the root cause within **5 business days**.
3. **Patch Development**: A fix will be developed in a private security advisory branch and validated against our automated test suite and security linters (`govulncheck`, `pre-commit`).
4. **Release & Credit**: A patched release will be tagged and deployed to `main`. We will publicly credit the reporter (unless anonymity is requested).
5. **Coordinated Disclosure**: Public details and CVE disclosures will be published after downstream consumers have had reasonable opportunity to update.

---

## Security Practices in IDXStock

IDXStock enforces several automated security and supply-chain controls:
- **Automated Dependency Updates**: Continuous dependency scanning and PRs via Dependabot.
- **Static Analysis & SAST**: Go code scanning with `govulncheck`, `go vet`, and CodeQL workflows.
- **Safe Database Interactions**: Parameterized queries via `pgx/v5` preventing SQL injection vulnerabilities.
- **Memory Safety & Error Boundaries**: Panic recovery middleware configured in the Fiber v3 HTTP engine.
