# Management API contract

`management.openapi.json` is a reviewed OpenAPI 3.0.3 subset used to generate the internal HTTP client and models. It contains seven operations for the initial FOCUS push connection lifecycle.

The source baseline is backend commit `923e81f496bb1be4fa5e064afce746d0c19ae755`. The new `GetIntegrationConnection` operation is supplied by the separate `feat/terraform-integration-contract` backend change. Stable operation identifiers for the other integration operations come from that change too. At SDK creation, the backend change was not deployed.

The full Swagger artifact generated from that backend checkout with its pinned `swag` v1.16.6 has SHA256 `487a650ffcb4d01de620c31586131d7067481ee0604efcae7fdbbd09ef7d11e7`. It is supporting research evidence, not a build input. The relevant source files are `controllers/integrations/integrations.go`, `controllers/integrations/registry.go`, `internal/user/user.go`, `responses/success.go`, and `responses/error.go` in the backend repository.

The reviewed differences from the full Swagger representation are deliberate.

- `/api/v1/me` has a session/API-key union in the full API description. This SDK authenticates with an API key, so its `Identity` schema follows `user.CurrentAPIKeyResponse`, including optional organization metadata and permission scopes.
- The list route's session-oriented tenant header is not required here. API-key authentication establishes the authorized tenant, and this client does not support session authentication.
- Connection models retain the common identity, provider, tenant, display name, status, ingestion mode, and timestamp fields. Telemetry, billing history, recommendation, and sandbox details are outside this SDK slice. The server can return those extra fields without changing this model.
- Required and nullable fields follow the Go JSON response and handler behavior. String IDs stay strings; `id` is an internal UUID and `connection_id` is the public `fpc_` identifier.
- FOCUS push creation requires a billing-account identifier. Provider, display name, and tenant are optional. The key's tenant is the default; omitting the provider allows the first billing push to set it.
- Rename and disconnect return success envelopes without a data payload. The facade also handles `204 No Content` without attempting JSON decoding.
- Error response schemas document the API shape, but the public facade exposes only the HTTP status and operation name. Arbitrary server text and raw bodies are not diagnostics.

Run `make generate` to regenerate with oapi-codegen v2.9.0. Run `make check-generated` to verify the committed client matches this snapshot. Neither command downloads the live CostGraph schema.
