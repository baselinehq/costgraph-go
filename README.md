# CostGraph Go management SDK

Use this Go module to manage FOCUS push connections through the CostGraph API. It also reads API-key identity, tenant-scoped integration connections, and the FOCUS billing-provider catalog.

This is an initial management SDK, not complete API coverage. The single-connection read endpoint requires the accompanying backend Terraform integration contract change. Existing hosted or self-hosted deployments may not contain that change yet.

## Configure the client

Pass the CostGraph API deployment URL and API key explicitly. The SDK does not load environment variables or saved CLI credentials and has no hosted default or fallback.

The URL is the deployment root, such as `https://api.costgraph.ai`, or a reverse-proxy prefix, such as `https://costgraph.example.com/backend`. Do not append `/api/v1`. Every SDK operation appends its API path under the supplied prefix. The URL must not contain user credentials, a query, or a fragment.

```go
client, err := costgraph.New(os.Getenv("COSTGRAPH_API_URL"), os.Getenv("COSTGRAPH_API_KEY"))
if err != nil {
    return err
}

identity, err := client.CurrentIdentity(ctx)
```

Import the package with `github.com/baselinehq/costgraph-go`. The client uses `X-API-Key`, a 30-second request timeout, and standard certificate verification. HTTP is supported for local test deployments. Supply a configured HTTP client with `costgraph.WithHTTPClient` when needed. The SDK copies that client and disables redirects so credentials do not follow a redirect to another endpoint.

The same client constructor supports self-hosted API roots and reverse proxies.

```go
selfHosted, err := costgraph.New("https://costgraph.internal.example.com", apiKey)
proxied, err := costgraph.New("https://platform.example.com/costgraph", apiKey)
```

Those clients call `/api/v1/...` and `/costgraph/api/v1/...` respectively. Use the API host rather than the dashboard host.

API keys are tenant-bound. Connection listing uses that authorized tenant. Creating a push connection can omit `TenantID` to use the key's tenant; providing a different tenant does not grant access to it. Identity organization fields are optional pointers because the API can omit them.

## Manage a push connection

The following calls create a billing account connection and read it back by its internal UUID.

```go
connection, err := client.CreatePushConnection(ctx, costgraph.CreatePushConnectionRequest{
    Provider:          "aws",
    ProviderAccountID: "billing-account-123",
    DisplayName:       "Production billing",
})
if err != nil {
    return err
}

connection, err = client.GetConnection(ctx, connection.ID)
```

`connection.ID` is the internal UUID used for lookup, rename, and deletion. `connection.ConnectionID` is the public `fpc_` identifier used when sending billing records. The SDK does not ingest those records in this initial slice.

Use `RenameConnection(ctx, id, name)` to update the display name. Provider, billing-account ID, and tenant changes are not update operations in this SDK. `DeleteConnection(ctx, id)` disconnects the integration and stops collection through that connection. It is a destructive action.

`ListIntegrationProviders` reads the FOCUS billing-provider catalog. It does not list every kind of CostGraph integration or establish that each provider has a supported managed Terraform resource.

## Handle API errors

HTTP failures return `*costgraph.APIError`, with `StatusCode` and `Operation`. Its string contains only a safe operation name and HTTP status. It does not include arbitrary server messages, response bodies, credentials, or request URLs. Inspect the status with `errors.As`. A 404 establishes absence; authentication, permission, network, and server failures do not.

Cancellation remains detectable with `errors.Is(err, context.Canceled)`. Mutation methods accept successful JSON envelopes and `204 No Content`. The SDK does not retry writes automatically.

## Build and generate

Go 1.26 or later is required. The public API contract is reviewed and committed at `api/management.openapi.json`. It describes the subset exposed here, including API-key authentication and typed response envelopes. Generated HTTP and model code lives under `internal/api`; the public facade owns configuration and safe error handling.

```sh
make build
make lint
make check-generated
```

Run `make generate` after an intentional contract change. It uses oapi-codegen v2.9.0. Ordinary builds do not download a live API schema or regenerate files.

Build, vet, and generated-diff checks establish local compilation and reproducibility. Real resource lifecycle verification belongs to the Terraform provider acceptance scenario against an isolated CostGraph backend.
