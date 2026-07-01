# Migration Plan (Draft)

## From CLI to Platform
This document outlines the step-by-step migration from the existing CLI tool to the full AppSec platform.

## Phase 0 (Current)
- Keep root main.go unchanged
- Scaffold new packages

## Phase 1
- Implement storage interfaces
- Add multi-DB support (connection strings per tenant)

## Phase 2
- Add auth layer (JWT + API keys)
- Web UI prototype

## Phase 3
- Full platform release
- API endpoints
- Evidence storage (local FS)

## Revert Strategy
To revert any milestone, use:
  git revert <commit_hash>
  git push origin feature/refactor/core
