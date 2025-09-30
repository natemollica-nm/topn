package utils

import (
	"fmt"
	"strings"
)

func ParseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, fmt.Errorf("empty size")
	}
	
	unit := int64(1)
	switch {
	case strings.HasSuffix(s, "TB"), strings.HasSuffix(s, "T"):
		unit = 1 << 40
		s = strings.TrimSuffix(strings.TrimSuffix(s, "TB"), "T")
	case strings.HasSuffix(s, "GB"), strings.HasSuffix(s, "G"):
		unit = 1 << 30
		s = strings.TrimSuffix(strings.TrimSuffix(s, "GB"), "G")
	case strings.HasSuffix(s, "MB"), strings.HasSuffix(s, "M"):
		unit = 1 << 20
		s = strings.TrimSuffix(strings.TrimSuffix(s, "MB"), "M")
	case strings.HasSuffix(s, "KB"), strings.HasSuffix(s, "K"):
		unit = 1 << 10
		s = strings.TrimSuffix(strings.TrimSuffix(s, "KB"), "K")
	}
	
	s = strings.TrimSpace(s)
	var val float64
	if _, err := fmt.Sscanf(s, "%f", &val); err != nil {
		return 0, fmt.Errorf("invalid size: %q", s)
	}
	return int64(val * float64(unit)), nil
}

func HumanSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%dB", n)
	}
	const (
		k = 1 << 10
		m = 1 << 20
		g = 1 << 30
		t = 1 << 40
	)
	switch {
	case n >= t:
		return fmt.Sprintf("%.1fT", float64(n)/float64(t))
	case n >= g:
		return fmt.Sprintf("%.1fG", float64(n)/float64(g))
	case n >= m:
		return fmt.Sprintf("%.1fM", float64(n)/float64(m))
	default:
		return fmt.Sprintf("%.1fK", float64(n)/float64(k))
	}
}

type MultiFlag []string

func (m *MultiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *MultiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}