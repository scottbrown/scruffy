package scruffy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	dryRun      bool
	prefix      string
	target      string
	description string
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean IP Access rules",
	Long:  `Clean IP Access rules from Cloudflare. Can clean all, by prefix, by target, or by description.`,
}

var cleanAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Clean all IP Access rules",
	Long:  `Clean all IP Access rules from the specified zone.`,
	RunE:  runCleanAll,
}

var cleanPrefixCmd = &cobra.Command{
	Use:   "prefix [PREFIX]",
	Short: "Clean IP Access rules with a specific prefix",
	Long:  `Clean all IP Access rules that start with the specified prefix.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCleanPrefix,
}

var cleanTargetCmd = &cobra.Command{
	Use:   "target [TARGET]",
	Short: "Clean a specific IP Access rule",
	Long:  `Clean a specific IP Access rule by its target (IP, CIDR, or ASN).`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCleanTarget,
}

var cleanDescriptionCmd = &cobra.Command{
	Use:   "description [DESCRIPTION]",
	Short: "Clean IP Access rules by description",
	Long:  `Clean all IP Access rules that contain the specified description/note.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCleanDescription,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.AddCommand(cleanAllCmd)
	cleanCmd.AddCommand(cleanPrefixCmd)
	cleanCmd.AddCommand(cleanTargetCmd)
	cleanCmd.AddCommand(cleanDescriptionCmd)

	cleanCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be deleted without actually deleting")
}

func runCleanAll(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	
	_, client, err := setupClient(ctx)
	if err != nil {
		return err
	}

	rules, err := client.ListAccessRules(ctx)
	if err != nil {
		return fmt.Errorf("failed to list access rules: %w", err)
	}

	return deleteRules(ctx, client, rules, "all rules")
}

func runCleanPrefix(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	prefix := args[0]
	
	_, client, err := setupClient(ctx)
	if err != nil {
		return err
	}

	rules, err := client.ListAccessRules(ctx)
	if err != nil {
		return fmt.Errorf("failed to list access rules: %w", err)
	}

	filteredRules := client.FilterRulesByPrefix(rules, prefix)
	return deleteRules(ctx, client, filteredRules, fmt.Sprintf("rules with prefix %q", prefix))
}

func runCleanTarget(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	target := args[0]
	
	_, client, err := setupClient(ctx)
	if err != nil {
		return err
	}

	rules, err := client.ListAccessRules(ctx)
	if err != nil {
		return fmt.Errorf("failed to list access rules: %w", err)
	}

	filteredRules := client.FilterRulesByTarget(rules, target)
	return deleteRules(ctx, client, filteredRules, fmt.Sprintf("rules targeting %q", target))
}

func runCleanDescription(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	description := args[0]
	
	_, client, err := setupClient(ctx)
	if err != nil {
		return err
	}

	rules, err := client.ListAccessRules(ctx)
	if err != nil {
		return fmt.Errorf("failed to list access rules: %w", err)
	}

	filteredRules := client.FilterRulesByDescription(rules, description)
	return deleteRules(ctx, client, filteredRules, fmt.Sprintf("rules containing description %q", description))
}

func setupClient(ctx context.Context) (string, *Client, error) {
	resolvedZoneID := zoneID
	
	if zoneID == "" && zoneName != "" {
		tempClient, err := NewClient(apiToken, "")
		if err != nil {
			return "", nil, fmt.Errorf("failed to create temporary client: %w", err)
		}
		
		resolvedZoneID, err = tempClient.ResolveZoneID(ctx, zoneName)
		if err != nil {
			return "", nil, fmt.Errorf("failed to resolve zone ID: %w", err)
		}
	}

	client, err := NewClient(apiToken, resolvedZoneID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create client: %w", err)
	}

	return resolvedZoneID, client, nil
}

func deleteRules(ctx context.Context, client *Client, rules []AccessRule, description string) error {
	if len(rules) == 0 {
		fmt.Printf("No %s found.\n", description)
		return nil
	}

	fmt.Printf("Found %d %s:\n", len(rules), description)
	for _, rule := range rules {
		fmt.Printf("  - %s (%s) - %s\n", rule.Target, rule.Mode, rule.Notes)
	}

	if dryRun {
		fmt.Printf("\nDry run mode: would delete %d rules\n", len(rules))
		return nil
	}

	fmt.Printf("\nDeleting %d rules...\n", len(rules))
	var errors []error
	deleted := 0

	for _, rule := range rules {
		if err := client.DeleteAccessRule(ctx, rule); err != nil {
			fmt.Printf("Failed to delete rule %s: %v\n", rule.Target, err)
			errors = append(errors, err)
		} else {
			fmt.Printf("Deleted: %s\n", rule.Target)
			deleted++
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to delete %d out of %d rules", len(errors), len(rules))
	}

	fmt.Printf("\nSuccessfully deleted %d rules\n", deleted)
	return nil
}