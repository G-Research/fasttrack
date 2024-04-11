package oidc

import (
	"context"
	"slices"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rotisserie/eris"

	"github.com/G-Research/fasttrackml/pkg/common/config/auth"
)

// ClientProvider provides an interface to work with OIDC provider.
type ClientProvider interface {
	// Verify makes Access Token verification.
	Verify(ctx context.Context, accessToken string) (*User, error)
	// IsAdmin makes check that current user is Admin user.
	IsAdmin(user *User) bool
}

// Client represents OIDC client.
type Client struct {
	config   *auth.Config
	verifier *oidc.IDTokenVerifier
}

// NewClient creates new OIDC client
func NewClient(ctx context.Context, config *auth.Config) (*Client, error) {
	provider, err := oidc.NewProvider(ctx, config.AuthOIDCProviderEndpoint)
	if err != nil {
		return nil, eris.Wrap(err, "error creating OIDC provider")
	}

	return &Client{
		config:   config,
		verifier: provider.Verifier(&oidc.Config{ClientID: config.AuthOIDCClientID}),
	}, nil
}

// Verify makes Access Token verification.
func (c Client) Verify(ctx context.Context, accessToken string) (*User, error) {
	idToken, err := c.verifier.Verify(ctx, accessToken)
	if err != nil {
		return nil, eris.Wrap(err, "error verifying access token")
	}
	// Extract custom claims.
	var claims struct {
		Groups []string `json:"groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, eris.Wrap(err, "error extracting token claims")
	}
	return &User{Groups: claims.Groups}, nil
}

// IsAdmin makes check that current user is Admin user.
func (c Client) IsAdmin(user *User) bool {
	return slices.Contains(user.Groups, c.config.AuthOIDCAdminRole)
}
