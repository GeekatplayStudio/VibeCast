# AgenticSFU Workspace Rules & Super Agent Guard Protocols

## Core Directive for All Agents

Every AI agent working on **AgenticSFU** MUST strictly adhere to the following quality, security, and architectural rules:

### 1. Roadmap Alignment & Non-Drift Guard
- Before adding any new feature, inspect `master_roadmap_analysis.md` and `feature_matrix.md`.
- Verify that new code aligns with planned architecture and does not drift from core principles.
- Maintain clear justification and documentation for every file edit or new component.

### 2. Comprehensive Documentation Integrity
- Whenever code changes are made, immediately update `README.md`, `walkthrough.md`, and inline comments.
- Maintain full docstrings on all exported Go functions, structs, and interfaces.

### 3. Continuous Security Auditing
- Enforce strict input validation on all HTTP endpoints (REST API, WHIP, MCP JSON-RPC, WebSocket).
- Verify HMAC-SHA256 JWT access token signatures for room joins.
- Check dependencies for vulnerabilities and use safe DTLS/SRTP WebRTC practices.

### 4. Automated Testing & Verification
- No task is complete until automated tests are created or updated, and both `go test ./pkg/... ./test/...` and `npm run build` pass with zero errors.

### 5. Multi-Language (i18n) Support Standard
- All user-facing UI text strings MUST be registered in `ui/src/i18n/translations.ts` supporting English, Spanish, French, German, Japanese, and Chinese.
