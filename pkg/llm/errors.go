package llm

import (
	"errors"
	"fmt"

	openaisdk "github.com/openai/openai-go/v3"
)

// HTTPError preserves the provider's HTTP status while retaining the original error.
type HTTPError struct {
	StatusCode int
	Err        error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("provider returned %d: %v", e.StatusCode, e.Err)
}

func (e *HTTPError) Unwrap() error { return e.Err }

// NormalizeHTTPError exposes HTTP status errors consistently across OpenAI-compatible clients.
func NormalizeHTTPError(err error) error {
	if err == nil {
		return nil
	}

	var normalized *HTTPError
	if errors.As(err, &normalized) {
		return err
	}

	var providerErr *openaisdk.Error
	if errors.As(err, &providerErr) {
		return &HTTPError{StatusCode: providerErr.StatusCode, Err: err}
	}

	return err
}
