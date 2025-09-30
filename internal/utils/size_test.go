package utils

import "testing"

func TestParseSize(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		hasError bool
	}{
		{"1G", 1 << 30, false},
		{"1.5GB", int64(1.5 * (1 << 30)), false},
		{"500M", 500 << 20, false},
		{"1024K", 1 << 20, false},
		{"1T", 1 << 40, false},
		{"", 0, true},
		{"invalid", 0, true},
	}
	
	for _, tt := range tests {
		got, err := ParseSize(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseSize(%q) expected error", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseSize(%q) error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("ParseSize(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		}
	}
}

func TestHumanSize(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{512, "512B"},
		{1024, "1.0K"},
		{1536, "1.5K"},
		{1048576, "1.0M"},
		{1073741824, "1.0G"},
		{1099511627776, "1.0T"},
	}
	
	for _, tt := range tests {
		got := HumanSize(tt.input)
		if got != tt.expected {
			t.Errorf("HumanSize(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}