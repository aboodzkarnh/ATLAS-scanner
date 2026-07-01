# Atlas Scanner Architecture (Draft)

## Overview
Atlas Scanner is evolving from a CLI vulnerability scanner into an AppSec platform.

## Current State (Milestone 0)
- CLI tool in root main.go (backward compatible)
- New packages scaffolded: core, scanner, storage, integrations, evidence, reports

## Future Milestones
1. Milestone 1: Storage layer (multi-DB per tenant)
2. Milestone 2: Auth layer (JWT + API keys)
3. Milestone 3: Web UI + API
4. Milestone 4: OOB integration (Interactsh self-hosted)

## Technology Stack
- Go 1.22
- SQLite (MVP), PostgreSQL (multi-tenant)
- GORM (ORM)
- Interactsh (OOB)
