package llm

import (
	"errors"
	"testing"

	openaisdk "github.com/openai/openai-go/v3"
)

func TestNormalizeHTTPError(t *testing.T) {
	err := NormalizeHTTPError(&openaisdk.Error{StatusCode: 429})

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != 429 {
		t.Fatalf("error = %#v", err)
	}
}
