package bootstrap

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"backend/internal/infrastructure/database/postgres"
)

// discardLogger returns a *slog.Logger that discards all output.
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// ---------------------------------------------------------------------------
// loadConfig
// ---------------------------------------------------------------------------

func TestLoadConfig_DefaultPort(t *testing.T) {
	t.Setenv("PORT", "")
	cfg := loadConfig()
	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}
}

func TestLoadConfig_ExplicitPort(t *testing.T) {
	t.Setenv("PORT", "9090")
	cfg := loadConfig()
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
}

func TestLoadConfig_DefaultSchema(t *testing.T) {
	t.Setenv("BLUEPRINT_DB_SCHEMA", "")
	cfg := loadConfig()
	if cfg.DB.Schema != "public" {
		t.Errorf("expected default schema %q, got %q", "public", cfg.DB.Schema)
	}
}

func TestLoadConfig_ExplicitSchema(t *testing.T) {
	t.Setenv("BLUEPRINT_DB_SCHEMA", "myschema")
	cfg := loadConfig()
	if cfg.DB.Schema != "myschema" {
		t.Errorf("expected schema %q, got %q", "myschema", cfg.DB.Schema)
	}
}

func TestLoadConfig_DefaultSSLMode(t *testing.T) {
	t.Setenv("BLUEPRINT_DB_SSLMODE", "")
	cfg := loadConfig()
	if cfg.DB.SSLMode != "disable" {
		t.Errorf("expected default sslmode %q, got %q", "disable", cfg.DB.SSLMode)
	}
}

func TestLoadConfig_ExplicitSSLMode(t *testing.T) {
	t.Setenv("BLUEPRINT_DB_SSLMODE", "require")
	cfg := loadConfig()
	if cfg.DB.SSLMode != "require" {
		t.Errorf("expected sslmode %q, got %q", "require", cfg.DB.SSLMode)
	}
}

func TestLoadConfig_SentryDSN_Set(t *testing.T) {
	const dsn = "https://key@o0.ingest.sentry.io/0"
	t.Setenv("SENTRY_DSN", dsn)
	cfg := loadConfig()
	if cfg.SentryDSN != dsn {
		t.Errorf("expected SentryDSN %q, got %q", dsn, cfg.SentryDSN)
	}
}

func TestLoadConfig_SentryDSN_Unset(t *testing.T) {
	t.Setenv("SENTRY_DSN", "")
	cfg := loadConfig()
	if cfg.SentryDSN != "" {
		t.Errorf("expected empty SentryDSN, got %q", cfg.SentryDSN)
	}
}

// ---------------------------------------------------------------------------
// validateConfig
// ---------------------------------------------------------------------------

func TestValidateConfig_AllPresent(t *testing.T) {
	cfg := Config{
		Port: 8080,
		DB:   dbConfigFull(),
	}
	if err := validateConfig(cfg, discardLogger()); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestValidateConfig_AllMissing(t *testing.T) {
	cfg := Config{} // all DB fields are zero-value (empty strings)
	err := validateConfig(cfg, discardLogger())
	if err == nil {
		t.Fatal("expected a ConfigError, got nil")
	}

	var cfgErr *ConfigError
	if !errors.As(err, &cfgErr) {
		t.Fatalf("expected *ConfigError, got %T: %v", err, err)
	}
	if len(cfgErr.Issues) != 5 {
		t.Errorf("expected 5 issues (one per required field), got %d: %v", len(cfgErr.Issues), cfgErr.Issues)
	}
}

func TestValidateConfig_MissingIndividualFields(t *testing.T) {
	base := dbConfigFull()

	tests := []struct {
		name        string
		mutate      func(cfg *Config)
		wantIssues  int
		issueSubstr string
	}{
		{
			name:        "missing host",
			mutate:      func(c *Config) { c.DB.Host = "" },
			wantIssues:  1,
			issueSubstr: "BLUEPRINT_DB_HOST",
		},
		{
			name:        "missing port",
			mutate:      func(c *Config) { c.DB.Port = "" },
			wantIssues:  1,
			issueSubstr: "BLUEPRINT_DB_PORT",
		},
		{
			name:        "missing database",
			mutate:      func(c *Config) { c.DB.Database = "" },
			wantIssues:  1,
			issueSubstr: "BLUEPRINT_DB_DATABASE",
		},
		{
			name:        "missing username",
			mutate:      func(c *Config) { c.DB.Username = "" },
			wantIssues:  1,
			issueSubstr: "BLUEPRINT_DB_USERNAME",
		},
		{
			name:        "missing password",
			mutate:      func(c *Config) { c.DB.Password = "" },
			wantIssues:  1,
			issueSubstr: "BLUEPRINT_DB_PASSWORD",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{Port: 8080, DB: base}
			tc.mutate(&cfg)

			err := validateConfig(cfg, discardLogger())
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var cfgErr *ConfigError
			if !errors.As(err, &cfgErr) {
				t.Fatalf("expected *ConfigError, got %T", err)
			}
			if len(cfgErr.Issues) != tc.wantIssues {
				t.Errorf("expected %d issue(s), got %d: %v", tc.wantIssues, len(cfgErr.Issues), cfgErr.Issues)
			}

			found := false
			for _, issue := range cfgErr.Issues {
				if contains(issue, tc.issueSubstr) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected an issue mentioning %q, got: %v", tc.issueSubstr, cfgErr.Issues)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// jitteredBackoff
// ---------------------------------------------------------------------------

func TestJitteredBackoff_NonNegative(t *testing.T) {
	for attempt := 1; attempt <= 10; attempt++ {
		d := jitteredBackoff(attempt)
		if d < 0 {
			t.Errorf("attempt %d: jitteredBackoff returned negative duration: %v", attempt, d)
		}
	}
}

func TestJitteredBackoff_NeverExceedsMaxDelay(t *testing.T) {
	for attempt := 1; attempt <= 20; attempt++ {
		d := jitteredBackoff(attempt)
		if d > maxDelay {
			t.Errorf("attempt %d: jitteredBackoff returned %v which exceeds maxDelay %v", attempt, d, maxDelay)
		}
	}
}

// ---------------------------------------------------------------------------
// probeWithRetry
// ---------------------------------------------------------------------------

// mockPinger is a local test double implementing Pinger.
type mockPinger struct {
	results []error // results[i] is returned on the (i+1)-th call; last value is repeated
	calls   int
}

func (m *mockPinger) PingContext(_ context.Context) error {
	if m.calls >= len(m.results) {
		return m.results[len(m.results)-1]
	}
	err := m.results[m.calls]
	m.calls++
	return err
}

func TestProbeWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	p := &mockPinger{results: []error{nil}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := probeWithRetry(ctx, "test-service", p, discardLogger())
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestProbeWithRetry_SuccessAfterRetries(t *testing.T) {
	sentinel := errors.New("not ready")
	p := &mockPinger{results: []error{sentinel, sentinel, nil}}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := probeWithRetry(ctx, "test-service", p, discardLogger())
	if err != nil {
		t.Errorf("expected nil after eventual success, got: %v", err)
	}
}

func TestProbeWithRetry_FailsAfterAllAttempts(t *testing.T) {
	sentinel := errors.New("always failing")
	// Always return an error so all maxAttempts are exhausted.
	p := &mockPinger{results: []error{sentinel}}

	// Use a deadline that is short enough that the jittered sleeps between
	// retries get cancelled, preventing the test from taking O(minutes).
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := probeWithRetry(ctx, "test-service", p, discardLogger())
	if err == nil {
		t.Fatal("expected error after all attempts failed, got nil")
	}
}

func TestProbeWithRetry_ContextCancelledMidRetry(t *testing.T) {
	sentinel := errors.New("unavailable")
	p := &mockPinger{results: []error{sentinel}}

	ctx, cancel := context.WithCancel(context.Background())
	// Cancel immediately so the first inter-retry sleep is interrupted.
	cancel()

	err := probeWithRetry(ctx, "test-service", p, discardLogger())
	if err == nil {
		t.Fatal("expected error when context is cancelled, got nil")
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// dbConfigFull returns a postgres.DBConfig with all required fields populated.
func dbConfigFull() postgres.DBConfig {
	return postgres.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		Database: "testdb",
		Username: "user",
		Password: "secret",
		Schema:   "public",
		SSLMode:  "disable",
	}
}

// contains reports whether s contains substr.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
