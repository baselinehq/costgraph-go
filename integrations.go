package costgraph

import (
	"context"
	"errors"

	"github.com/baselinehq/costgraph-go/internal/api"
	"github.com/google/uuid"
)

func (c *Client) CurrentIdentity(ctx context.Context) (*Identity, error) {
	response, err := c.api.CurrentIdentity(ctx)
	result, err := decode[api.IdentityResponse]("read current identity", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.SubjectType == "" {
		return nil, errors.New("read current identity: missing identity in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) ListIntegrationProviders(ctx context.Context) ([]IntegrationProvider, error) {
	response, err := c.api.ListFocusProviders(ctx)
	result, err := decode[api.ProviderListResponse]("list integration providers", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data == nil {
		return nil, errors.New("list integration providers: missing provider list in CostGraph API response")
	}
	return result.Data, nil
}

func (c *Client) ListConnections(ctx context.Context) ([]Connection, error) {
	response, err := c.api.ListIntegrationConnections(ctx)
	result, err := decode[api.ConnectionListResponse]("list connections", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data == nil {
		return nil, errors.New("list connections: missing connection list in CostGraph API response")
	}
	return result.Data, nil
}

func (c *Client) GetConnection(ctx context.Context, id string) (*Connection, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.GetIntegrationConnection(ctx, id)
	result, err := decode[api.ConnectionResponse]("read connection", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("read connection: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) CreatePushConnection(ctx context.Context, request CreatePushConnectionRequest) (*Connection, error) {
	response, err := c.api.CreateFocusPushConnection(ctx, request)
	result, err := decode[api.ConnectionResponse]("create FOCUS push connection", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("create FOCUS push connection: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) RenameConnection(ctx context.Context, id, name string) error {
	if err := requireConnectionID(id); err != nil {
		return err
	}
	response, err := c.api.RenameIntegrationConnection(ctx, id, api.RenameConnectionRequest{DisplayName: name})
	_, err = decode[api.SuccessResponse]("rename connection", response, err)
	return err
}

func (c *Client) DeleteConnection(ctx context.Context, id string) error {
	if err := requireConnectionID(id); err != nil {
		return err
	}
	response, err := c.api.DisconnectIntegrationConnection(ctx, id)
	_, err = decode[api.SuccessResponse]("disconnect connection", response, err)
	return err
}

func requireConnectionID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("connection id must be the internal connection UUID, not its fpc_ connection_id")
	}
	return nil
}
