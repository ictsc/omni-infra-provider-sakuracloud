// Package main runs the Sakura Cloud Omni infrastructure provider.
package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/siderolabs/omni/client/pkg/client"
	"github.com/siderolabs/omni/client/pkg/infra"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ictsc/omni-infra-provider-sakuracloud/internal/provider"
	providermeta "github.com/ictsc/omni-infra-provider-sakuracloud/internal/provider/meta"
	"github.com/ictsc/omni-infra-provider-sakuracloud/internal/sakura"
)

//go:embed data/schema.json
var schema string

var flags struct {
	omniEndpoint       string
	serviceAccountKey  string
	providerName       string
	allowMutations     bool
	insecureSkipVerify bool
}

var rootCmd = &cobra.Command{
	Use:          "omni-infra-provider-sakuracloud",
	Short:        "Sakura Cloud infrastructure provider for Omni",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if flags.omniEndpoint == "" {
			return fmt.Errorf("--omni-api-endpoint is required")
		}

		logger, err := zap.NewProduction()
		if err != nil {
			return fmt.Errorf("create logger: %w", err)
		}
		defer logger.Sync() //nolint:errcheck

		provisioner := provider.NewProvisioner(sakura.NewGuardedClient(flags.allowMutations))
		infraProvider, err := infra.NewProvider(providermeta.ProviderID, provisioner, infra.ProviderConfig{
			Name:        flags.providerName,
			Description: "Dynamic Talos workers on Sakura Cloud",
			Schema:      schema,
		})
		if err != nil {
			return fmt.Errorf("create Omni infra provider: %w", err)
		}

		clientOptions := []client.Option{client.WithInsecureSkipTLSVerify(flags.insecureSkipVerify)}
		if flags.serviceAccountKey != "" {
			clientOptions = append(clientOptions, client.WithServiceAccount(flags.serviceAccountKey))
		}

		return infraProvider.Run(
			cmd.Context(),
			logger,
			infra.WithOmniEndpoint(flags.omniEndpoint),
			infra.WithClientOptions(clientOptions...),
			infra.WithEncodeRequestIDsIntoTokens(),
		)
	},
}

func init() {
	rootCmd.Flags().StringVar(&flags.omniEndpoint, "omni-api-endpoint", os.Getenv("OMNI_ENDPOINT"), "Omni API endpoint")
	rootCmd.Flags().StringVar(&flags.serviceAccountKey, "omni-service-account-key", os.Getenv("OMNI_SERVICE_ACCOUNT_KEY"), "Omni infrastructure provider key")
	rootCmd.Flags().StringVar(&providermeta.ProviderID, "id", providermeta.ProviderID, "Omni infrastructure provider ID")
	rootCmd.Flags().StringVar(&flags.providerName, "provider-name", "Sakura Cloud", "provider name shown in Omni")
	rootCmd.Flags().BoolVar(&flags.allowMutations, "allow-mutations", false, "allow Sakura Cloud resource mutations")
	rootCmd.Flags().BoolVar(&flags.insecureSkipVerify, "insecure-skip-verify", false, "skip Omni TLS certificate verification")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
