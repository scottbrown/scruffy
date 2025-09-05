package scruffy

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiToken string
	zoneID   string
	zoneName string
)

var rootCmd = &cobra.Command{
	Use:   "scruffy",
	Short: "Clean Cloudflare IP Access rules",
	Long: `Scruffy is a CLI tool for cleaning Cloudflare IP Access rules.
It can clean all records, records with specific prefixes, or specific records.
Records can be IP addresses, CIDR blocks, or ASNs.`,
  Version: Version(),
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiToken, "token", "", "Cloudflare API token (discouraged, use CLOUDFLARE_API_TOKEN env var instead)")
	rootCmd.PersistentFlags().StringVar(&zoneID, "zone-id", "", "Cloudflare Zone ID")
	rootCmd.PersistentFlags().StringVar(&zoneName, "zone-name", "", "Cloudflare Zone name (alternative to zone-id)")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if apiToken == "" {
			apiToken = os.Getenv("CLOUDFLARE_API_TOKEN")
		}
		if apiToken == "" {
			return fmt.Errorf("API token is required. Set CLOUDFLARE_API_TOKEN environment variable or use --token flag (not recommended)")
		}

		if zoneID == "" && zoneName == "" {
			return fmt.Errorf("either --zone-id or --zone-name must be specified")
		}

		return nil
	}
}
