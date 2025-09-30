package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/natemollica/topn/internal/scanner"
)

type scanProgressMsg struct {
	progress float64
	current  string
}

func (m Model) startScanWithProgress() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start progress updates
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			
			progress := 0.0
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					progress += 0.02
					if progress > 0.9 {
						progress = 0.1
					}
					// Send progress update (this won't work in real implementation)
					// In practice, you'd need a more sophisticated progress tracking
				}
			}
		}()

		s := scanner.New(m.config)
		results, stats := s.Scan()
		return scanCompleteMsg{results: results, stats: stats}
	})
}