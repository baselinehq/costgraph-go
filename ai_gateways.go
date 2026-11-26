package costgraph

import (
	"context"
	"errors"

	"github.com/baselinehq/costgraph-go/internal/api"
)

func (c *Client) CreateAIGatewayConnection(ctx context.Context, request AIGatewayConnectRequest) (*Connection, error) {
	response, err := c.api.CreateAIGatewayConnection(ctx, request)
	result, err := decode[api.ConnectionResponse]("create AI gateway connection", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("create AI gateway connection: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) GetAIGatewaySettings(ctx context.Context, id string) (*AIGatewaySettings, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.GetAIGatewaySettings(ctx, id)
	result, err := decode[api.AIGatewaySettingsResponse]("read AI gateway settings", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.IngestMode == "" {
		return nil, errors.New("read AI gateway settings: missing settings in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) UpdateAIGatewayConnection(ctx context.Context, id string, request AIGatewayUpdateRequest) (*Connection, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.UpdateAIGatewayConnection(ctx, id, request)
	result, err := decode[api.ConnectionResponse]("update AI gateway connection", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("update AI gateway connection: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) SearchCatalogModels(ctx context.Context, query string) ([]CatalogModel, error) {
	response, err := c.api.SearchAIGatewayModels(ctx, &api.SearchAIGatewayModelsParams{Q: &query})
	result, err := decode[api.CatalogModelListResponse]("search catalog models", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data == nil {
		return nil, errors.New("search catalog models: missing model list in CostGraph API response")
	}
	return result.Data, nil
}
