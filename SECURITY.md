# Security Policy

## Supported Versions

This project is currently in active development. Security updates are typically only applied to the latest version.

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < latest | :x:                |

## Reporting a Vulnerability

**IMPORTANT**: Do not report security vulnerabilities through public GitHub issues.

Please report security vulnerabilities by:

1. Using GitHub's private vulnerability reporting feature at https://github.com/chrisw-dev/golang-mock-oauth2-server/security/advisories/new
2. Alternatively, emailing `security.golang-mock-oauth2-server@beardycoder.co.uk`

Please include the following information in your report:

- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

This information will help us triage your report more quickly.

## Policy

- We will acknowledge receipt of your vulnerability report within 3 business days
- We will provide an initial assessment of the report within 10 business days
- We will keep you informed of our progress towards resolving the issue
- We ask that you do not publicly disclose the issue until we have had a chance to address it

## Important Notice

This is a mock OAuth2 server intended for development and testing purposes only. It does not implement proper security measures required for real OAuth2 authentication and should never be used in production environments.

When reporting security issues, please keep in mind the intended use case of this software.