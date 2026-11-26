package costgraph

import (
	"context"
	"errors"

	"github.com/baselinehq/costgraph-go/internal/api"
	"github.com/google/uuid"
)

func (c *Client) CreateKubernetesCluster(ctx context.Context, organizationID string, request CreateKubernetesClusterRequest) (*KubernetesCluster, error) {
	if err := requireTenantIDs(organizationID, request.TenantID); err != nil {
		return nil, err
	}
	response, err := c.api.CreateKubernetesCluster(ctx, organizationID, request)
	result, err := decode[api.KubernetesClusterResponse]("reserve Kubernetes cluster", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("reserve Kubernetes cluster: missing registration in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) GetKubernetesCluster(ctx context.Context, organizationID, id string) (*KubernetesCluster, error) {
	if err := requireKubernetesClusterIDs(organizationID, id); err != nil {
		return nil, err
	}
	response, err := c.api.GetKubernetesCluster(ctx, organizationID, id)
	result, err := decode[api.KubernetesClusterResponse]("read Kubernetes cluster", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("read Kubernetes cluster: missing registration in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) DeleteKubernetesCluster(ctx context.Context, organizationID, id string) error {
	if err := requireKubernetesClusterIDs(organizationID, id); err != nil {
		return err
	}
	response, err := c.api.DeleteKubernetesCluster(ctx, organizationID, id)
	_, err = decode[api.SuccessResponse]("disconnect Kubernetes cluster", response, err)
	return err
}

func requireKubernetesClusterIDs(organizationID, id string) error {
	if err := requireTenantIDs(organizationID); err != nil {
		return err
	}
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("cluster id must be a UUID")
	}
	return nil
}
