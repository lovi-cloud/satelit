package sshkey

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/whywaita/satelit/pkg/api"

	isucon_sshkey "github.com/whywaita/isucon-sshkey"
	"go.uber.org/zap"
)

// Endpoint for dev
const (
	DevEndpoint = "dev"
)

// check mode
var (
	IsDev = false
)

// NewClient create a client of isucon_sshkey and create mock server if dev.
func NewClient(endpoint, hmacSecretKey string, logger *zap.Logger) (*isucon_sshkey.Client, error) {
	if DevEndpoint == endpoint {
		IsDev = true
		return NewMockClient(logger)
	}

	client, err := isucon_sshkey.NewClient(endpoint, hmacSecretKey, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create ISUCON portal client: %w", err)
	}

	return client, nil
}

// NewMockClient create mock client for dev.
func NewMockClient(logger *zap.Logger) (*isucon_sshkey.Client, error) {
	logger.Debug("detect dev sshkey endpoint. running for mock server")

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, api.AdminKeys)
		})

	ts := httptest.NewServer(mux)

	return isucon_sshkey.NewClient(ts.URL, "dev", logger)
}
