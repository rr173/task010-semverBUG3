package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Regression: /api/parse must return prerelease and build as empty JSON
// arrays ([]) when the version has no prerelease or build metadata,
// not as JSON null.

func TestParseReturnsEmptyArrayNotNull(t *testing.T) {
	api := New()
	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/parse", "application/json",
		strings.NewReader(`{"version":"1.0.0"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	ver, ok := body["version"].(map[string]any)
	if !ok {
		t.Fatal("missing version field")
	}

	// prerelease must be a non-nil array
	pre, ok := ver["prerelease"]
	if !ok {
		t.Fatal("missing prerelease field")
	}
	preArr, ok := pre.([]any)
	if !ok {
		t.Errorf("prerelease should be array, got %T (value: %v)", pre, pre)
	} else if preArr == nil {
		t.Error("prerelease array is nil, expected empty []")
	}

	// build must be a non-nil array
	bld, ok := ver["build"]
	if !ok {
		t.Fatal("missing build field")
	}
	bldArr, ok := bld.([]any)
	if !ok {
		t.Errorf("build should be array, got %T (value: %v)", bld, bld)
	} else if bldArr == nil {
		t.Error("build array is nil, expected empty []")
	}
}
