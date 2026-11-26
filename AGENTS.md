# CostGraph Go management SDK

## Scope

- This module is the importable management client used by the CostGraph Terraform provider. Keep Terraform schemas, state, diagnostics, imports, and lifecycle decisions in the provider repository.
- The client covers API-key identity, provider and connection discovery, FOCUS push lifecycle, organization tenant administration, AWS role connections, Google Cloud impersonation connections, and AI gateway settings. It is not a complete CostGraph SDK.
- Match the actual public API. Inspect its routes, authentication, request/response envelopes, tenant boundaries, and permissions before adding an operation. A dashboard endpoint does not establish API-key access.
- The read and management contracts are introduced by separate backend Terraform integration and management changes. Do not claim availability on deployments that lack those changes.
- Organization tenant and cloud administration require an active organization owner and explicit `organization:manage` scope. Do not weaken or implicitly add that scope. Cloud SDK requests support AWS roles and Google Cloud impersonation only.

## Contract and client

- `api/management.openapi.json` is a reviewed OpenAPI 3 subset, not a live schema download. Keep its provenance accurate. Add stable operation identifiers and typed success/error envelopes.
- Generate `internal/api/client.gen.go` with `make generate`. The generator version is pinned in the Makefile. Commit generated changes; never hand-edit generated Go or fetch a changing upstream schema during a normal build.
- Keep the public facade in package `costgraph`. Preserve null/omitted values where the API distinguishes them. Use optional scalar and map pointers for partial updates; nil preserves, while an explicit empty value clears only when the API supports it. Add only the API capabilities the caller needs.
- Every operation uses the configured deployment URL and `X-API-Key`. Support deployment roots and reverse-proxy path prefixes. Never fall back to the hosted API or read a browser/CLI credential file.
- Retain normal certificate verification. Do not follow redirects with credentials. Do not expose API keys, request URLs, raw response bodies, or arbitrary server messages in diagnostics.
- Use contexts and bounded request timeouts. Do not retry non-idempotent writes automatically. Treat permission/transport errors separately from resource absence.
- Connection `id` is the internal UUID. `connection_id` is the public `fpc_` ingestion identifier. Do not interchange them.

## Verification and shipping

- `make build` compiles the module; `make lint` runs Go vet; `make generate` regenerates the client; `make check-generated` detects generated drift. `make fmt` formats Go.
- No unnecessary tests. Do not add mocked HTTP/CRUD tests or generated boilerplate as substitutes for API lifecycle proof. The provider's real acceptance scenario must prove create, API readback, an empty follow-up plan, update, import, destroy, and API-confirmed absence against an isolated real backend.
- Record checks accurately. Compilation and generation are not evidence that a live API works. Skipped acceptance runs are not passing coverage.
- Keep changes small and close API gaps before provider resources. Do not create releases, remote repositories, or unsupported compatibility claims without agreed publication details.

## Durable notes

Save durable notes, ideas, research, audits, and plans in `~/notes`, using Obsidian-compatible Markdown. Begin with YAML `date`, `topic`, and a quoted ISO 8601 `timestamp` with local timezone offset. Preserve creation metadata and add `updated` when revising. Include Context, Objectives, and Outcomes, distinguish verified work from proposals, and keep credentials out. Link related notes and store supporting artifacts in the topic's `assets/` directory. Keep code and publishable docs in their repositories.

The shared discovery note is `~/notes/CostGraph/2026-11-26 CostGraph Terraform provider discovery.md`. Pass its path and these conventions to agents and include its clickable link in handoffs. The shared template is `~/notes/Templates/Agent note.md`.
