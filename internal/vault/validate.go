package vault

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationOptions controls which validation rules are applied to secrets.
type ValidationOptions struct {
	Enabled        bool
	RequireNonEmpty bool
	ForbiddenKeys  []string
	MaxValueLength int
}

// DefaultValidationOptions returns sensible defaults (validation disabled).
func DefaultValidationOptions() ValidationOptions {
	return ValidationOptions{
		Enabled:        false,
		RequireNonEmpty: false,
		ForbiddenKeys:  nil,
		MaxValueLength: 0,
	}
}

// ValidationError holds all violations found during validation.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d violation(s): %s",
		len(e.Violations), strings.Join(e.Violations, "; "))
}

// ValidateSecrets checks the provided secrets map against the given options.
// Returns a *ValidationError if any violations are found, nil otherwise.
// If opts.Enabled is false the call is a no-op.
func ValidateSecrets(secrets map[string]map[string]string, opts ValidationOptions) error {
	if !opts.Enabled {
		return nil
	}

	var violations []string

	for path, kv := range secrets {
		for key, value := range kv {
			if opts.RequireNonEmpty && strings.TrimSpace(value) == "" {
				violations = append(violations,
					fmt.Sprintf("%s/%s: value must not be empty", path, key))
			}

			if opts.MaxValueLength > 0 && len(value) > opts.MaxValueLength {
				violations = append(violations,
					fmt.Sprintf("%s/%s: value length %d exceeds max %d",
						path, key, len(value), opts.MaxValueLength))
			}

			for _, forbidden := range opts.ForbiddenKeys {
				if strings.EqualFold(key, forbidden) {
					violations = append(violations,
						fmt.Sprintf("%s/%s: key is forbidden", path, key))
				}
			}
		}
	}

	if len(violations) > 0 {
		return &ValidationError{Violations: violations}
	}
	return nil
}

// IsValidationError reports whether err is a *ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
