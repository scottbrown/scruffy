package scruffy

import (
	"testing"
)

func TestFilterRulesByPrefix(t *testing.T) {
	client := &Client{}

	rules := []AccessRule{
		{ID: "1", Target: "192.168.1.1", Mode: "block"},
		{ID: "2", Target: "192.168.2.1", Mode: "block"},
		{ID: "3", Target: "10.0.0.1", Mode: "block"},
		{ID: "4", Target: "AS64496", Mode: "block"},
	}

	tests := []struct {
		name     string
		prefix   string
		expected int
	}{
		{
			name:     "filter by 192.168 prefix",
			prefix:   "192.168",
			expected: 2,
		},
		{
			name:     "filter by 10. prefix",
			prefix:   "10.",
			expected: 1,
		},
		{
			name:     "filter by AS prefix",
			prefix:   "AS",
			expected: 1,
		},
		{
			name:     "no matches",
			prefix:   "172.16",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.FilterRulesByPrefix(rules, tt.prefix)
			if len(result) != tt.expected {
				t.Errorf("expected %d rules, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestFilterRulesByTarget(t *testing.T) {
	client := &Client{}

	rules := []AccessRule{
		{ID: "1", Target: "192.168.1.1", Mode: "block"},
		{ID: "2", Target: "192.168.2.1", Mode: "block"},
		{ID: "3", Target: "192.168.1.1", Mode: "allow"},
	}

	result := client.FilterRulesByTarget(rules, "192.168.1.1")
	if len(result) != 2 {
		t.Errorf("expected 2 rules matching target, got %d", len(result))
	}

	result = client.FilterRulesByTarget(rules, "nonexistent")
	if len(result) != 0 {
		t.Errorf("expected 0 rules for nonexistent target, got %d", len(result))
	}
}

func TestFilterRulesByDescription(t *testing.T) {
	client := &Client{}

	rules := []AccessRule{
		{ID: "1", Target: "192.168.1.1", Notes: "temporary block for testing"},
		{ID: "2", Target: "192.168.2.1", Notes: "permanent security block"},
		{ID: "3", Target: "10.0.0.1", Notes: "temp access rule"},
		{ID: "4", Target: "AS64496", Notes: ""},
	}

	tests := []struct {
		name        string
		description string
		expected    int
	}{
		{
			name:        "filter by 'temp' description",
			description: "temp",
			expected:    2,
		},
		{
			name:        "filter by 'block' description",
			description: "block",
			expected:    2,
		},
		{
			name:        "filter by 'security' description",
			description: "security",
			expected:    1,
		},
		{
			name:        "no matches",
			description: "nonexistent",
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.FilterRulesByDescription(rules, tt.description)
			if len(result) != tt.expected {
				t.Errorf("expected %d rules, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		apiToken  string
		zoneID    string
		expectErr bool
	}{
		{
			name:      "valid token and zone ID",
			apiToken:  "test-token",
			zoneID:    "test-zone-id",
			expectErr: false,
		},
		{
			name:      "empty token",
			apiToken:  "",
			zoneID:    "test-zone-id",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiToken, tt.zoneID)
			
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				if client != nil {
					t.Error("expected nil client on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if client == nil {
					t.Error("expected non-nil client")
				}
				if client.zoneID != tt.zoneID {
					t.Errorf("expected zone ID %s, got %s", tt.zoneID, client.zoneID)
				}
			}
		})
	}
}

func TestAccessRuleStruct(t *testing.T) {
	rule := AccessRule{
		ID:     "test-id",
		Target: "192.168.1.1",
		Mode:   "block",
		Notes:  "test note",
		Scope:  "zone",
	}

	if rule.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got %s", rule.ID)
	}
	if rule.Target != "192.168.1.1" {
		t.Errorf("expected Target '192.168.1.1', got %s", rule.Target)
	}
	if rule.Mode != "block" {
		t.Errorf("expected Mode 'block', got %s", rule.Mode)
	}
	if rule.Notes != "test note" {
		t.Errorf("expected Notes 'test note', got %s", rule.Notes)
	}
	if rule.Scope != "zone" {
		t.Errorf("expected Scope 'zone', got %v", rule.Scope)
	}
}