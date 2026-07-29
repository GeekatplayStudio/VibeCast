---
name: super_agent_auditor
description: Audits code additions against roadmap requirements, checks for security vulnerabilities, verifies tests and lints, and ensures documentation and multi-language support are in sync.
---

# Super Agent Auditor Skill

This skill defines the autonomous auditing process for checking all additions, commits, and pull requests in **AgenticSFU**.

## Audit Workflow

### 1. Roadmap & Architecture Audit
- Check `master_roadmap_analysis.md` and `feature_matrix.md`.
- Verify that every change has an explicit architectural explanation.
- Ensure no accidental feature drift or scope creep occurs.

### 2. Security & Vulnerability Audit
- Verify HMAC-SHA256 JWT signature verification logic in `pkg/auth/token.go`.
- Validate input sanitization on WebSockets, REST APIs, and WHIP endpoints.
- Check CORS headers and ensure no secret tokens are leaked in logs.
- Audit dependencies in `go.mod` and `ui/package.json`.

### 3. Test & Lint Verification
- Run `go test ./pkg/... ./test/...` to ensure 100% test passing.
- Run TypeScript build check `npm run build` inside `ui/`.

### 4. Documentation & Commit Sync
- Ensure `README.md`, `walkthrough.md`, and code comments are fully updated.
- Verify user-facing UI text strings are localized in `ui/src/i18n/translations.ts`.
