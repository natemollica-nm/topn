package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/natemollica-nm/topn/internal/scanner"
	"github.com/natemollica-nm/topn/internal/ui"
	"github.com/natemollica-nm/topn/internal/utils"
)

var version = "dev"

func main() {
	var (
		dir      string
		minStr   string
		topN     int
		workers  int
		exclVals utils.MultiFlag
		remove   bool
		tui      bool
		showVer  bool
	)

	flag.StringVar(&dir, "dir", os.Getenv("HOME"), "root directory to scan")
	flag.StringVar(&minStr, "min", "1G", "minimum file size (e.g. 1G, 500M, 250K)")
	flag.IntVar(&topN, "top", 50, "keep only top N largest files")
	flag.IntVar(&workers, "workers", 0, "number of workers (default: 4*GOMAXPROCS)")
	flag.Var(&exclVals, "exclude", "glob/path to exclude (repeatable)")
	flag.BoolVar(&remove, "remove", false, "interactively remove files")
	flag.BoolVar(&tui, "tui", false, "use interactive terminal UI")
	flag.BoolVar(&showVer, "version", false, "show version")
	flag.Parse()

	if showVer {
		fmt.Printf("topn %s\n", version)
		return
	}

	if workers <= 0 {
		workers = 4 * runtime.GOMAXPROCS(0)
	}

	minBytes, err := utils.ParseSize(minStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing size: %v\n", err)
		os.Exit(1)
	}

	root, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving directory: %v\n", err)
		os.Exit(1)
	}

	if st, err := os.Stat(root); err != nil || !st.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: '%s' is not a valid directory\n", root)
		os.Exit(1)
	}

	config := scanner.Config{
		Root:     root,
		MinBytes: minBytes,
		TopN:     topN,
		Workers:  workers,
		Excludes: exclVals,
	}

	// Use TUI if requested or if remove flag is set
	if tui || remove {
		model := ui.NewModel(config)
		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)
		
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Classic CLI mode with enhanced output
	fmt.Printf("🔍 Scanning %s for files >= %s...\n", root, minStr)
	
	s := scanner.New(config)
	start := time.Now()
	results, stats := s.Scan()
	elapsed := time.Since(start).Round(time.Millisecond)

	fmt.Printf("\n✅ Scan complete in %s\n", elapsed)
	fmt.Printf("📊 Files seen: %d, kept: %d (>= %s)\n\n", 
		stats.FilesSeen, stats.FilesKept, minStr)

	if len(results) == 0 {
		fmt.Println("🎉 No large files found!")
		return
	}

	printResults(results)
	
	if len(results) > 0 {
		fmt.Printf("\n💡 Tip: Use -tui or -remove for interactive file management\n")
	}
}

func printResults(results []scanner.FileItem) {
	fmt.Printf("%-5s %-10s %s\n", "Rank", "Size", "Path")
	fmt.Printf("%-5s %-10s %s\n", "----", "----", strings.Repeat("-", 50))
	
	for i, item := range results {
		rank := fmt.Sprintf("#%d", i+1)
		size := utils.HumanSize(item.Size)
		
		// Truncate long paths for better display
		path := item.Path
		if len(path) > 70 {
			path = "..." + path[len(path)-67:]
		}
		
		fmt.Printf("%-5s %-10s %s\n", rank, size, path)
	}
}