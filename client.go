package scruffy

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

type Client struct {
	api    *cloudflare.API
	zoneID string
}

type AccessRule struct {
	ID     string
	Target string
	Mode   string
	Notes  string
	Scope  string
}

func NewClient(apiToken string, zoneID string) (*Client, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudflare client: %w", err)
	}

	return &Client{
		api:    api,
		zoneID: zoneID,
	}, nil
}

func (c *Client) ResolveZoneID(ctx context.Context, zoneName string) (string, error) {
	zones, err := c.api.ListZones(ctx, zoneName)
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %w", err)
	}

	for _, zone := range zones {
		if zone.Name == zoneName {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("zone %q not found", zoneName)
}

func (c *Client) ListAccessRules(ctx context.Context) ([]AccessRule, error) {
	var allRules []AccessRule
	page := 1

	for {
		zoneRules, err := c.api.ListZoneAccessRules(ctx, c.zoneID, cloudflare.AccessRule{}, page)
		if err != nil {
			return nil, fmt.Errorf("failed to list zone access rules on page %d: %w", page, err)
		}

		// Add rules from this page
		for _, rule := range zoneRules.Result {
			allRules = append(allRules, AccessRule{
				ID:     rule.ID,
				Target: rule.Configuration.Value,
				Mode:   rule.Mode,
				Notes:  rule.Notes,
				Scope:  "zone",
			})
		}

		// Check if we have more pages
		if len(zoneRules.Result) == 0 || page >= zoneRules.TotalPages {
			break
		}

		page++
	}

	return allRules, nil
}

func (c *Client) DeleteAccessRule(ctx context.Context, rule AccessRule) error {
	switch rule.Scope {
	case "zone":
		_, err := c.api.DeleteZoneAccessRule(ctx, c.zoneID, rule.ID)
		return err
	case "account":
		_, err := c.api.DeleteAccountAccessRule(ctx, "", rule.ID)
		return err
	default:
		return fmt.Errorf("unsupported access rule scope: %v", rule.Scope)
	}
}

func (c *Client) FilterRulesByPrefix(rules []AccessRule, prefix string) []AccessRule {
	var filtered []AccessRule
	for _, rule := range rules {
		if strings.HasPrefix(rule.Target, prefix) {
			filtered = append(filtered, rule)
		}
	}
	return filtered
}

func (c *Client) FilterRulesByTarget(rules []AccessRule, target string) []AccessRule {
	var filtered []AccessRule
	for _, rule := range rules {
		if rule.Target == target {
			filtered = append(filtered, rule)
		}
	}
	return filtered
}

func (c *Client) FilterRulesByDescription(rules []AccessRule, description string) []AccessRule {
	var filtered []AccessRule
	for _, rule := range rules {
		if strings.Contains(rule.Notes, description) {
			filtered = append(filtered, rule)
		}
	}
	return filtered
}