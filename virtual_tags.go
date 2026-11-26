package costgraph

import (
	"context"
	"errors"
	"strings"

	"github.com/baselinehq/costgraph-go/internal/api"
	"github.com/google/uuid"
)

const ScopeVirtualTagsRead = "virtual-tags:read"
const ScopeVirtualTagsWrite = "virtual-tags:write"

func virtualTagTenant(tenantID string) error {
	if _, err := uuid.Parse(tenantID); err != nil {
		return errors.New("virtual tag tenant_id must be a UUID")
	}
	return nil
}

func virtualTagIdentity(tenantID, id string) error {
	if err := virtualTagTenant(tenantID); err != nil {
		return err
	}
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("virtual tag id must be a key UUID")
	}
	return nil
}

func virtualTagResult(operation string, response api.VirtualTagResponse, err error) (*VirtualTag, error) {
	if err != nil {
		return nil, err
	}
	if response.Data.ID == "" || response.Data.TenantID == "" || response.Data.Revision == "" || response.Data.Rules == nil {
		return nil, errors.New(operation + ": incomplete virtual tag configuration in CostGraph API response")
	}
	return &response.Data, nil
}

func (c *Client) GetVirtualTag(ctx context.Context, tenantID, id string) (*VirtualTag, error) {
	if err := virtualTagIdentity(tenantID, id); err != nil {
		return nil, err
	}
	response, err := c.api.GetVirtualTagConfiguration(ctx, id, &api.GetVirtualTagConfigurationParams{TenantID: tenantID})
	result, err := decode[api.VirtualTagResponse]("read virtual tag", response, err)
	return virtualTagResult("read virtual tag", result, err)
}

func (c *Client) LookupVirtualTag(ctx context.Context, tenantID, key string) (*VirtualTag, error) {
	if err := virtualTagTenant(tenantID); err != nil {
		return nil, err
	}
	response, err := c.api.LookupVirtualTagConfiguration(ctx, &api.LookupVirtualTagConfigurationParams{TenantID: tenantID, Key: key})
	result, err := decode[api.VirtualTagResponse]("look up virtual tag", response, err)
	return virtualTagResult("look up virtual tag", result, err)
}

func (c *Client) ListVirtualTags(ctx context.Context, tenantID string) ([]VirtualTagMetadata, error) {
	if err := virtualTagTenant(tenantID); err != nil {
		return nil, err
	}
	items := []VirtualTagMetadata{}
	params := &api.ListVirtualTagConfigurationsParams{TenantID: tenantID}
	seen := map[string]bool{}
	for {
		response, err := c.api.ListVirtualTagConfigurations(ctx, params)
		result, err := decode[api.VirtualTagPageResponse]("list virtual tags", response, err)
		if err != nil {
			return nil, err
		}
		if result.Data.Items == nil {
			return nil, errors.New("list virtual tags: missing metadata list in CostGraph API response")
		}
		items = append(items, result.Data.Items...)
		cursor := result.Data.NextCursor
		if cursor == nil || *cursor == "" {
			return items, nil
		}
		if seen[*cursor] {
			return nil, errors.New("list virtual tags: repeated pagination cursor")
		}
		seen[*cursor] = true
		params.After = cursor
	}
}

func (c *Client) CreateVirtualTag(ctx context.Context, tenantID string, request VirtualTagRequest) (*VirtualTag, error) {
	if err := virtualTagTenant(tenantID); err != nil {
		return nil, err
	}
	response, err := c.api.CreateVirtualTagConfiguration(ctx, &api.CreateVirtualTagConfigurationParams{TenantID: tenantID}, request)
	result, err := decode[api.VirtualTagResponse]("create virtual tag", response, err)
	return virtualTagResult("create virtual tag", result, err)
}

func (c *Client) UpdateVirtualTag(ctx context.Context, tenantID, id string, request VirtualTagRequest) (*VirtualTag, error) {
	if err := virtualTagIdentity(tenantID, id); err != nil {
		return nil, err
	}
	response, err := c.api.UpdateVirtualTagConfiguration(ctx, id, &api.UpdateVirtualTagConfigurationParams{TenantID: tenantID}, request)
	result, err := decode[api.VirtualTagResponse]("update virtual tag", response, err)
	return virtualTagResult("update virtual tag", result, err)
}

func (c *Client) DeleteVirtualTag(ctx context.Context, tenantID, id, revision string) error {
	if err := virtualTagIdentity(tenantID, id); err != nil {
		return err
	}
	if revision == "" || strings.ContainsAny(revision, "\r\n\"") {
		return errors.New("virtual tag revision must be an unquoted aggregate revision")
	}
	response, err := c.api.DeleteVirtualTagConfiguration(ctx, id, &api.DeleteVirtualTagConfigurationParams{TenantID: tenantID, Revision: `"` + revision + `"`})
	_, err = decode[api.SuccessResponse]("delete virtual tag", response, err)
	return err
}
