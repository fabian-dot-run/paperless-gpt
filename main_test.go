package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDocument containing extra parameters for testing
type TestDocument struct {
	ID         int
	Title      string
	Tags       []string
	FailUpdate bool // simulate update failure
}

// Use this for TestCases in your tests
type TestCase struct {
	name           string
	documents      []TestDocument
	expectedCount  int
	expectedError  string
	updateResponse int // HTTP status code for update response
}

// Test our HTTP-Client
func TestCreateCustomHTTPClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify custom header
		assert.Equal(t, "paperless-gpt", r.Header.Get("X-Title"), "Expected X-Title header")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Get custom client
	client := createCustomHTTPClient()
	require.NotNil(t, client, "HTTP client should not be nil")

	// Make a request
	resp, err := client.Get(server.URL)
	require.NoError(t, err, "Request should not fail")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK response")
}

func TestGetTimeoutFromEnv(t *testing.T) {
	const envVar = "TEST_TIMEOUT_SECONDS"

	tests := []struct {
		name           string
		value          string
		setEnv         bool
		defaultSeconds int
		expected       time.Duration
	}{
		{name: "unset uses default", setEnv: false, defaultSeconds: 30, expected: 30 * time.Second},
		{name: "empty uses default", value: "", setEnv: true, defaultSeconds: 30, expected: 30 * time.Second},
		{name: "valid value parsed", value: "120", setEnv: true, defaultSeconds: 30, expected: 120 * time.Second},
		{name: "zero means no timeout", value: "0", setEnv: true, defaultSeconds: 30, expected: 0},
		{name: "invalid falls back to default", value: "abc", setEnv: true, defaultSeconds: 30, expected: 30 * time.Second},
		{name: "negative falls back to default", value: "-5", setEnv: true, defaultSeconds: 30, expected: 30 * time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(envVar, tc.value)
			} else {
				os.Unsetenv(envVar)
			}
			assert.Equal(t, tc.expected, getTimeoutFromEnv(envVar, tc.defaultSeconds))
		})
	}
}
