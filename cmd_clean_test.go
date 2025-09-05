package scruffy

import (
	"context"
	"testing"
)

func TestSetupClient(t *testing.T) {
	originalZoneID := zoneID
	originalZoneName := zoneName
	originalAPIToken := apiToken

	defer func() {
		zoneID = originalZoneID
		zoneName = originalZoneName
		apiToken = originalAPIToken
	}()

	tests := []struct {
		name      string
		zoneID    string
		zoneName  string
		apiToken  string
		expectErr bool
	}{
		{
			name:     "valid zone ID",
			zoneID:   "test-zone-id",
			apiToken: "test-token",
		},
		{
			name:      "empty token",
			zoneID:    "test-zone-id",
			apiToken:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zoneID = tt.zoneID
			zoneName = tt.zoneName
			apiToken = tt.apiToken

			ctx := context.Background()
			resolvedZoneID, client, err := setupClient(ctx)

			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if client == nil {
					t.Error("expected non-nil client")
				}
				if resolvedZoneID != tt.zoneID {
					t.Errorf("expected zone ID %s, got %s", tt.zoneID, resolvedZoneID)
				}
			}
		})
	}
}

func TestDeleteRulesFunction(t *testing.T) {
	ctx := context.Background()
	client := &Client{}
	
	rules := []AccessRule{
		{ID: "1", Target: "192.168.1.1", Mode: "block", Notes: "test"},
		{ID: "2", Target: "10.0.0.1", Mode: "allow", Notes: "allow rule"},
	}

	originalDryRun := dryRun
	defer func() { dryRun = originalDryRun }()

	t.Run("dry run mode", func(t *testing.T) {
		dryRun = true
		err := deleteRules(ctx, client, rules, "test rules")
		if err != nil {
			t.Errorf("dry run should not error: %v", err)
		}
	})

	t.Run("no rules", func(t *testing.T) {
		dryRun = false
		err := deleteRules(ctx, client, []AccessRule{}, "no rules")
		if err != nil {
			t.Errorf("no rules should not error: %v", err)
		}
	})
}

func TestCleanCommands(t *testing.T) {
	if cleanCmd.Use != "clean" {
		t.Errorf("expected Use 'clean', got %s", cleanCmd.Use)
	}

	if cleanAllCmd.Use != "all" {
		t.Errorf("expected Use 'all', got %s", cleanAllCmd.Use)
	}

	if cleanPrefixCmd.Use != "prefix [PREFIX]" {
		t.Errorf("expected Use 'prefix [PREFIX]', got %s", cleanPrefixCmd.Use)
	}

	if cleanTargetCmd.Use != "target [TARGET]" {
		t.Errorf("expected Use 'target [TARGET]', got %s", cleanTargetCmd.Use)
	}

	if cleanDescriptionCmd.Use != "description [DESCRIPTION]" {
		t.Errorf("expected Use 'description [DESCRIPTION]', got %s", cleanDescriptionCmd.Use)
	}

	dryRunFlag := cleanCmd.PersistentFlags().Lookup("dry-run")
	if dryRunFlag == nil {
		t.Error("expected --dry-run flag to exist")
	}
}