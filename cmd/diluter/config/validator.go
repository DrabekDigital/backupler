package config

import (
	"fmt"
	"math"

	"drabek.digital/cli-utils/backupler/cmd/diluter/definitions"
	"drabek.digital/cli-utils/backupler/cmd/diluter/helpers"
	"drabek.digital/cli-utils/backupler/cmd/diluter/strategies"
)

func ValidateConfig(c definitions.Config) error {

	// Check each policy
	for _, policy := range c.Policy {
		err := validatePolicy(policy)
		if err != nil {
			return err
		}
	}
	// Check interval <0, infinity) coverage
	err := validatePolicyCoverage(c)
	if err != nil {
		return err
	}
	return nil
}

func validatePolicy(p definitions.Policy) error {
	// From
	parsedFrom, err := helpers.ParseDays(p.From, false)
	if err != nil {
		return fmt.Errorf("`from` with %w", err)
	}
	// To
	parsedTo, err := helpers.ParseDays(p.To, true)
	if err != nil {
		return fmt.Errorf("`to` with %w", err)
	}

	// Interval
	if parsedTo < parsedFrom || parsedFrom == parsedTo {
		return fmt.Errorf("interval <`%d`, `%d`) is invalid", parsedFrom, parsedTo)
	}
	// Strategy
	if err = validateStrategy(p.Strategy); err != nil {
		return err
	}
	return nil
}

func validateStrategy(s definitions.Strategy) error {
	if s.Name == "" {
		return fmt.Errorf("strategy cannot be empty")
	}
	// Validate presence of window
	if s.Name != strategies.Delete && s.Name != strategies.Keep && s.Name != strategies.Dilute {
		return fmt.Errorf("strategy `%s` is unsupported", s.Name)
	}
	if (s.Name == strategies.Delete || s.Name == strategies.Keep) && s.Window != nil {
		return fmt.Errorf("strategy `%s` does not support window specified", s.Name)
	}
	if (s.Name == strategies.Dilute) && s.Window == nil {
		return fmt.Errorf("strategy `%s` has to have window specified", s.Name)
	}
	if (s.Name == strategies.Dilute) && s.Window != nil {
		_, err := helpers.ParseDays(*s.Window, false)
		if err != nil {
			return fmt.Errorf("`window` with %w", err)
		}
	}

	return nil
}

func validatePolicyCoverage(c definitions.Config) error {
	if len(c.Policy) == 0 {
		return fmt.Errorf("at least one policy has to be defined")
	}

	var startLine int64 = 0

	for _, policy := range c.Policy {
		from, err := helpers.ParseDays(policy.From, false)
		if err != nil {
			return err
		}
		to, err := helpers.ParseDays(policy.To, true)
		if err != nil {
			return err
		}
		if startLine != from {
			return fmt.Errorf("expected policy to start from `%d days`, instead `%s` given", startLine, policy.From)
		}
		startLine = to
	}
	endLine, err := helpers.ParseDays(c.Policy[len(c.Policy)-1].To, true)
	if err != nil {
		return err
	}
	if endLine != math.MaxInt64 {
		return fmt.Errorf("expected last policy to end with `infinity`, instead `%d` given", endLine)
	}

	return nil
}
