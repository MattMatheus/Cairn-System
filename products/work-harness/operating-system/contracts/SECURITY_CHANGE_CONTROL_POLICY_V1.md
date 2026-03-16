# Security Change Control Policy v1

Defines mandatory controls for deployment and security-affecting operating contract changes.

## Scope
- `dev -> prod` promotion path
- Any change to security-backing contracts listed in:
  - `products/work-harness/operating-system/contracts/SECURITY_BACKING_CONTRACT_PATHS.txt`

## Mandatory Launch Controls
1. Human operator confirmation is required for all `dev -> prod` transitions.
2. Dual confirmation phrase is required:
   - `CAIRN_SECURITY_CHANGE_ACK=LAUNCH AUTHORIZED`
3. One-time nonce is required for security policy changes:
   - generated per launch window
   - disclosed to operator
   - never reused after expiration
4. Security changes require dual OTP:
   - primary launch OTP (`CAIRN_SECURITY_CHANGE_OTP_PRIMARY`)
   - security policy OTP (`CAIRN_SECURITY_CHANGE_OTP_SECURITY`)
5. Nonce validity window is enforced by environment variable:
   - `CAIRN_SECURITY_GATE_WINDOW_MINUTES` (default `15`)

## Required Gate Variables
- `CAIRN_SECURITY_CHANGE_NONCE`
- `CAIRN_SECURITY_NONCE_ISSUED_AT_UTC` (ISO8601 UTC; example `2026-02-27T12:00:00Z`)
- `CAIRN_SECURITY_CHANGE_OTP_PRIMARY`
- `CAIRN_SECURITY_CHANGE_OTP_SECURITY`
- `CAIRN_SECURITY_CHANGE_ACK`
- Optional:
  - `CAIRN_SECURITY_GATE_WINDOW_MINUTES` (default `15`)
  - `CAIRN_SECURITY_GATE_BASE_REF` (default `origin/main`)

## CI Enforcement
- CI must run:
  - `products/work-harness/tools/check_security_nonce_gate.sh`
- If protected security contract files changed and gate inputs are missing/invalid, CI must fail.

## Operator Posture
- Friction for security downgrades is intentional.
- Any failed validation blocks launch and requires fresh nonce + OTP sequence.
