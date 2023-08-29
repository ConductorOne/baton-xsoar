package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	AccessToken string `mapstructure:"token"`
	Unsafe      bool   `mapstructure:"unsafe"`
	ApiUrl      string `mapstructure:"api-url"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.AccessToken == "" {
		return fmt.Errorf("an access token must be provided")
	}

	if cfg.ApiUrl == "" {
		return fmt.Errorf("the API URL of the Cortex XSOAR instance must be provided")
	}
	parsedApiUrl, err := url.Parse(cfg.ApiUrl)
	if err != nil {
		return fmt.Errorf("failed to parse the API URL: %w", err)
	}
	if parsedApiUrl.Scheme != "https" {
		return fmt.Errorf("the API URL must use the HTTPS scheme")
	}

	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("token", "", "Access token used to connect to the Cortex XSOAR API. ($BATON_TOKEN)")
	cmd.PersistentFlags().Bool("unsafe", false, "Allow insecure TLS connections to Cortex XSOAR instance. ($BATON_UNSAFE)")
	cmd.PersistentFlags().String("api-url", "", "The API URL of the Cortex XSOAR instance. ($BATON_API_URL)")
}
