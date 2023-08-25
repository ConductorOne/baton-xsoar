package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	AccessToken string `mapstructure:"token"`
	Unsafe      bool   `mapstructure:"unsafe"`
	Domain      string `mapstructure:"domain"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.AccessToken == "" {
		return fmt.Errorf("an access token must be provided")
	}

	if cfg.Domain == "" {
		return fmt.Errorf("a domain of the Cortex XSOAR instance must be provided")
	}

	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("token", "", "Access token used to connect to the Cortex XSOAR API. ($BATON_TOKEN)")
	cmd.PersistentFlags().Bool("unsafe", false, "Allow insecure TLS connections to Cortex XSOAR instance. ($BATON_UNSAFE)")
	cmd.PersistentFlags().String("domain", "", "Domain of the Cortex XSOAR instance. ($BATON_DOMAIN)")
}
