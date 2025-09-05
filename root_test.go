package scruffy

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&apiToken, "token", "", "test token")
	cmd.PersistentFlags().StringVar(&zoneID, "zone-id", "", "test zone id")
	cmd.PersistentFlags().StringVar(&zoneName, "zone-name", "", "test zone name")

	cmd.PersistentPreRunE = rootCmd.PersistentPreRunE

	tests := []struct {
		name      string
		envToken  string
		flagToken string
		zoneID    string
		zoneName  string
		expectErr bool
	}{
		{
			name:      "valid with env token and zone ID",
			envToken:  "env-token",
			zoneID:    "test-zone-id",
			expectErr: false,
		},
		{
			name:      "valid with flag token and zone name",
			flagToken: "flag-token",
			zoneName:  "example.com",
			expectErr: false,
		},
		{
			name:      "no token provided",
			expectErr: true,
		},
		{
			name:     "token but no zone",
			envToken: "test-token",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiToken = ""
			zoneID = ""
			zoneName = ""

			if tt.envToken != "" {
				os.Setenv("CLOUDFLARE_API_TOKEN", tt.envToken)
				defer os.Unsetenv("CLOUDFLARE_API_TOKEN")
			}

			if tt.flagToken != "" {
				apiToken = tt.flagToken
			}

			if tt.zoneID != "" {
				zoneID = tt.zoneID
			}

			if tt.zoneName != "" {
				zoneName = tt.zoneName
			}

			err := cmd.PersistentPreRunE(cmd, []string{})

			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestRootCommandFlags(t *testing.T) {
	if rootCmd.Use != "scruffy" {
		t.Errorf("expected Use 'scruffy', got %s", rootCmd.Use)
	}

	if rootCmd.Short != "Clean Cloudflare IP Access rules" {
		t.Errorf("expected Short description to match")
	}

	tokenFlag := rootCmd.PersistentFlags().Lookup("token")
	if tokenFlag == nil {
		t.Error("expected --token flag to exist")
	}

	zoneIDFlag := rootCmd.PersistentFlags().Lookup("zone-id")
	if zoneIDFlag == nil {
		t.Error("expected --zone-id flag to exist")
	}

	zoneNameFlag := rootCmd.PersistentFlags().Lookup("zone-name")
	if zoneNameFlag == nil {
		t.Error("expected --zone-name flag to exist")
	}
}