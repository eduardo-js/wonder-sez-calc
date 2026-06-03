package config

import (
	"testing"
	"time"
)

func TestLoad_Addr(t *testing.T) {
	tests := []struct {
		name     string
		envVal   string
		expected string
	}{
		{"default when empty", "", ":8080"},
		{"custom value", ":9090", ":9090"},
		{"custom host and port", "0.0.0.0:3000", "0.0.0.0:3000"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ADDR", tc.envVal)
			cfg := Load()
			if cfg.Addr != tc.expected {
				t.Errorf("Addr: got %q, want %q", cfg.Addr, tc.expected)
			}
		})
	}
}

func TestLoad_AllowedOrigins(t *testing.T) {
	type envEntry struct {
		set bool
		val string
	}
	defaultOrigins := []string{"http://localhost:5173"}

	tests := []struct {
		name     string
		env      envEntry
		expected []string
	}{
		{
			name:     "env unset uses default",
			env:      envEntry{set: false},
			expected: defaultOrigins,
		},
		{
			name:     "empty string uses default",
			env:      envEntry{set: true, val: ""},
			expected: defaultOrigins,
		},
		{
			name:     "only commas and spaces uses default",
			env:      envEntry{set: true, val: "  ,  ,  "},
			expected: defaultOrigins,
		},
		{
			name:     "single origin",
			env:      envEntry{set: true, val: "https://example.com"},
			expected: []string{"https://example.com"},
		},
		{
			name:     "multiple comma-separated",
			env:      envEntry{set: true, val: "https://a.com,https://b.com,https://c.com"},
			expected: []string{"https://a.com", "https://b.com", "https://c.com"},
		},
		{
			name:     "origins with surrounding spaces",
			env:      envEntry{set: true, val: "  https://a.com , https://b.com  "},
			expected: []string{"https://a.com", "https://b.com"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.env.set {
				t.Setenv("CORS_ALLOWED_ORIGINS", tc.env.val)
			}
			cfg := Load()
			if len(cfg.AllowedOrigins) != len(tc.expected) {
				t.Fatalf("AllowedOrigins len: got %d, want %d (got %v)", len(cfg.AllowedOrigins), len(tc.expected), cfg.AllowedOrigins)
			}
			for i, want := range tc.expected {
				if cfg.AllowedOrigins[i] != want {
					t.Errorf("AllowedOrigins[%d]: got %q, want %q", i, cfg.AllowedOrigins[i], want)
				}
			}
		})
	}
}

func TestLoad_Timeouts(t *testing.T) {
	cfg := Load()
	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("ReadTimeout: got %v, want 5s", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 10*time.Second {
		t.Errorf("WriteTimeout: got %v, want 10s", cfg.WriteTimeout)
	}
	if cfg.IdleTimeout != 120*time.Second {
		t.Errorf("IdleTimeout: got %v, want 120s", cfg.IdleTimeout)
	}
}
