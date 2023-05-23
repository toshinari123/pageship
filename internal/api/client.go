package api

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/oursky/pageship/internal/config"
	"github.com/oursky/pageship/internal/models"
)

type Client struct {
	endpoint string
	client   *http.Client
}

func NewClientWithTransport(endpoint string, transport http.RoundTripper) *Client {
	return &Client{
		endpoint: endpoint,
		client:   &http.Client{Transport: transport},
	}
}

func NewClient(endpoint string) *Client {
	return NewClientWithTransport(endpoint, http.DefaultTransport)
}

func (c *Client) CreateApp(ctx context.Context, appID string) (*APIApp, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps")
	if err != nil {
		return nil, err
	}

	req, err := newJSONRequest(ctx, "POST", endpoint, map[string]any{
		"id": appID,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[*APIApp](resp)
}

func (c *Client) ListApps(ctx context.Context) ([]APIApp, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[[]APIApp](resp)
}

func (c *Client) ConfigureApp(ctx context.Context, appID string, conf *config.AppConfig) (*APIApp, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps", appID, "config")
	if err != nil {
		return nil, err
	}

	req, err := newJSONRequest(ctx, "PUT", endpoint, map[string]any{
		"config": conf,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[*APIApp](resp)
}

func (c *Client) CreateSite(ctx context.Context, appID string, siteName string) (*APISite, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps", appID, "sites")
	if err != nil {
		return nil, err
	}

	req, err := newJSONRequest(ctx, "POST", endpoint, map[string]any{
		"name": siteName,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[*APISite](resp)
}

func (c *Client) ListSites(ctx context.Context, appID string) ([]APISite, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps", appID, "sites")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[[]APISite](resp)
}

func (c *Client) UpdateSite(
	ctx context.Context,
	appID string,
	siteName string,
	patch *SitePatchRequest,
) (*APISite, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps", appID, "sites", siteName)
	if err != nil {
		return nil, err
	}

	req, err := newJSONRequest(ctx, "PATCH", endpoint, patch)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[*APISite](resp)
}

func (c *Client) ListDeployments(ctx context.Context, appID string) ([]APIDeployment, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps", appID, "deployments")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[[]APIDeployment](resp)
}

func (c *Client) SetupDeployment(
	ctx context.Context,
	appID string,
	name string,
	files []models.FileEntry,
	siteConfig *config.SiteConfig,
) (*models.Deployment, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps", appID, "deployments")
	if err != nil {
		return nil, err
	}

	req, err := newJSONRequest(ctx, "POST", endpoint, map[string]any{
		"name":        name,
		"files":       files,
		"site_config": siteConfig,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[*models.Deployment](resp)
}

func (c *Client) UploadDeploymentTarball(
	ctx context.Context,
	appID string,
	deploymentID string,
	tarball io.Reader,
) (*models.Deployment, error) {
	endpoint, err := url.JoinPath(c.endpoint, "api", "v1", "apps", appID, "deployments", deploymentID, "tarball")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, tarball)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeJSONResponse[*models.Deployment](resp)
}
