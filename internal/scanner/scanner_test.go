package scanner

import (
	"testing"
)

func TestExcludes(t *testing.T) {
	ex := excludes{globs: []string{"*.log", "node_modules", "/tmp/*"}}
	
	tests := []struct {
		path     string
		expected bool
	}{
		{"/home/user/app.log", true},
		{"/home/user/data.txt", false},
		{"/home/user/node_modules/pkg", true},
		{"/tmp/file", true},
		{"/var/tmp/file", false},
	}
	
	for _, tt := range tests {
		if got := ex.match(tt.path); got != tt.expected {
			t.Errorf("match(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestKeepTopN(t *testing.T) {
	h := &minHeap{}
	
	// Add items
	keepTopN(h, FileItem{Size: 100, Path: "a"}, 3)
	keepTopN(h, FileItem{Size: 200, Path: "b"}, 3)
	keepTopN(h, FileItem{Size: 50, Path: "c"}, 3)
	keepTopN(h, FileItem{Size: 300, Path: "d"}, 3)
	
	if h.Len() != 3 {
		t.Errorf("heap size = %d, want 3", h.Len())
	}
	
	// Should contain 100, 200, 300 (not 50)
	min := (*h)[0].Size
	if min != 100 {
		t.Errorf("min size = %d, want 100", min)
	}
}