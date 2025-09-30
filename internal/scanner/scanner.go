package scanner

import (
	"container/heap"
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type FileItem struct {
	Size int64
	Path string
}

type Stats struct {
	FilesSeen int64
	FilesKept int64
}

type Config struct {
	Root     string
	MinBytes int64
	TopN     int
	Workers  int
	Excludes []string
}

type Scanner struct {
	config Config
	ex     excludes
}

type ProgressCallback func(current string, progress float64)

func New(config Config) *Scanner {
	return &Scanner{
		config: config,
		ex:     excludes{globs: config.Excludes},
	}
}

func (s *Scanner) Scan() ([]FileItem, Stats) {
	return s.ScanWithProgress(nil)
}

func (s *Scanner) ScanWithProgress(callback ProgressCallback) ([]FileItem, Stats) {
	ctx := context.Background()
	return s.ScanWithContext(ctx, callback)
}

func (s *Scanner) ScanWithContext(ctx context.Context, callback ProgressCallback) ([]FileItem, Stats) {
	var filesSeen, filesKept atomic.Int64
	
	h := &minHeap{}
	heap.Init(h)
	var mu sync.Mutex
	
	var wg sync.WaitGroup
	pathChan := make(chan string, 1000)
	
	// Start workers
	for i := 0; i < s.config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChan {
				select {
				case <-ctx.Done():
					return
				default:
				}
				
				filesSeen.Add(1)
				if callback != nil {
					callback(path, float64(filesSeen.Load())/1000.0) // Rough progress
				}
				
				if info, err := os.Lstat(path); err == nil && info.Mode().IsRegular() {
					if sz := info.Size(); sz >= s.config.MinBytes {
						mu.Lock()
						keepTopN(h, FileItem{Size: sz, Path: path}, s.config.TopN)
						mu.Unlock()
						filesKept.Add(1)
					}
				}
			}
		}()
	}
	
	// Walk filesystem
	go func() {
		defer close(pathChan)
		filepath.WalkDir(s.config.Root, func(path string, d os.DirEntry, err error) error {
			select {
			case <-ctx.Done():
				return filepath.SkipAll
			default:
			}
			
			if err != nil {
				return nil
			}
			if s.ex.match(path) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.Type()&os.ModeSymlink != 0 {
				return nil
			}
			if !d.IsDir() {
				select {
				case pathChan <- path:
				case <-ctx.Done():
					return filepath.SkipAll
				}
			}
			return nil
		})
	}()
	
	wg.Wait()

	// Extract results
	results := make([]FileItem, h.Len())
	for i := len(results) - 1; i >= 0; i-- {
		results[i] = heap.Pop(h).(FileItem)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Size > results[j].Size })

	return results, Stats{
		FilesSeen: filesSeen.Load(),
		FilesKept: filesKept.Load(),
	}
}

type minHeap []FileItem

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].Size < h[j].Size }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(FileItem)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func keepTopN(h *minHeap, it FileItem, n int) {
	if n <= 0 {
		return
	}
	if h.Len() < n {
		heap.Push(h, it)
	} else if (*h)[0].Size < it.Size {
		heap.Pop(h)
		heap.Push(h, it)
	}
}

type excludes struct{ globs []string }

func (e excludes) match(path string) bool {
	for _, g := range e.globs {
		if ok, _ := filepath.Match(g, path); ok {
			return true
		}
		if ok, _ := filepath.Match(g, filepath.Base(path)); ok {
			return true
		}
		if !strings.ContainsAny(g, "*?[]") && strings.Contains(path, g) {
			return true
		}
	}
	return false
}