package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/natemollica-nm/topn/internal/scanner"
	"github.com/natemollica-nm/topn/internal/utils"
)

type state int

const (
	stateScanning state = iota
	stateViewing
	stateConfirming
	stateHelp
)

type Model struct {
	state     state
	table     table.Model
	progress  progress.Model
	help      help.Model
	keys      keyMap
	results   []scanner.FileItem
	stats     scanner.Stats
	config    scanner.Config
	selected  map[int]bool
	message   string
	err       error
	width     int
	height    int
}

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Select    key.Binding
	Remove    key.Binding
	Rescan    key.Binding
	Help      key.Binding
	Quit      key.Binding
	Confirm   key.Binding
	Cancel    key.Binding
	SelectAll key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Remove, k.Rescan, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.SelectAll},
		{k.Remove, k.Rescan, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
	Down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
	Select:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "select/deselect")),
	Remove:    key.NewBinding(key.WithKeys("enter", "d"), key.WithHelp("enter/d", "delete selected")),
	Rescan:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rescan")),
	Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Confirm:   key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "confirm")),
	Cancel:    key.NewBinding(key.WithKeys("n", "esc"), key.WithHelp("n/esc", "cancel")),
	SelectAll: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "select all")),
}

type scanCompleteMsg struct {
	results []scanner.FileItem
	stats   scanner.Stats
}

type removeCompleteMsg struct {
	removed int
	errors  int
	message string
}

func NewModel(config scanner.Config) Model {
	columns := []table.Column{
		{Title: "Select", Width: 8},
		{Title: "Size", Width: 10},
		{Title: "Path", Width: 60},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = HeaderStyle.Copy().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)
	s.Selected = SelectedStyle.Copy()
	t.SetStyles(s)

	return Model{
		state:    stateScanning,
		table:    t,
		progress: progress.New(progress.WithDefaultGradient()),
		help:     help.New(),
		keys:     keys,
		config:   config,
		selected: make(map[int]bool),
	}
}

func (m Model) Init() tea.Cmd {
	return m.startScan()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.table.SetWidth(msg.Width - 4)
		m.table.SetHeight(msg.Height - 10)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateScanning:
			if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			}

		case stateViewing:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Help):
				m.state = stateHelp
				return m, nil
			case key.Matches(msg, m.keys.Select):
				if len(m.results) > 0 {
					idx := m.table.Cursor()
					m.selected[idx] = !m.selected[idx]
					m.updateTable()
				}
			case key.Matches(msg, m.keys.SelectAll):
				allSelected := len(m.selected) == len(m.results)
				m.selected = make(map[int]bool)
				if !allSelected {
					for i := range m.results {
						m.selected[i] = true
					}
				}
				m.updateTable()
			case key.Matches(msg, m.keys.Remove):
				if m.hasSelected() {
					m.state = stateConfirming
					return m, nil
				}
			case key.Matches(msg, m.keys.Rescan):
				m.state = stateScanning
				m.results = nil
				m.selected = make(map[int]bool)
				m.message = ""
				return m, m.startScan()
			}

		case stateConfirming:
			switch {
			case key.Matches(msg, m.keys.Confirm):
				m.state = stateScanning
				return m, m.removeSelected()
			case key.Matches(msg, m.keys.Cancel):
				m.state = stateViewing
				return m, nil
			}

		case stateHelp:
			if key.Matches(msg, m.keys.Help) || key.Matches(msg, m.keys.Quit) {
				m.state = stateViewing
				return m, nil
			}
		}

	case scanCompleteMsg:
		m.state = stateViewing
		m.results = msg.results
		m.stats = msg.stats
		m.updateTable()
		return m, nil

	case removeCompleteMsg:
		m.state = stateViewing
		m.message = msg.message
		m.selected = make(map[int]bool)
		return m, m.startScan()
	}

	if m.state == stateViewing {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateScanning:
		return m.scanningView()
	case stateViewing:
		return m.viewingView()
	case stateConfirming:
		return m.confirmingView()
	case stateHelp:
		return m.helpView()
	}
	return ""
}

func (m Model) scanningView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("🔍 TopN - Scanning Files"))
	b.WriteString("\n\n")
	b.WriteString(InfoStyle.Render(fmt.Sprintf("Scanning: %s", m.config.Root)))
	b.WriteString("\n\n")
	b.WriteString(ProgressStyle.Render("⠋ Finding large files..."))
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("Press q to quit"))
	return b.String()
}

func (m Model) viewingView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("🔍 TopN - Large File Scanner"))
	b.WriteString("\n\n")

	if len(m.results) > 0 {
		selectedCount := len(m.selected)
		b.WriteString(fmt.Sprintf(
			"Found %s files (%s kept >= %s) • %s selected\n\n",
			InfoStyle.Render(fmt.Sprintf("%d", m.stats.FilesSeen)),
			SuccessStyle.Render(fmt.Sprintf("%d", m.stats.FilesKept)),
			SizeStyle.Render(utils.HumanSize(m.config.MinBytes)),
			HeaderStyle.Render(fmt.Sprintf("%d", selectedCount)),
		))
		b.WriteString(m.table.View())
	} else {
		b.WriteString(InfoStyle.Render("No files found matching criteria"))
		b.WriteString("\n\n")
	}

	if m.message != "" {
		b.WriteString("\n")
		if strings.Contains(m.message, "error") {
			b.WriteString(ErrorStyle.Render(m.message))
		} else {
			b.WriteString(SuccessStyle.Render(m.message))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.help.View(m.keys))
	return b.String()
}

func (m Model) confirmingView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("⚠️  Confirm Deletion"))
	b.WriteString("\n\n")

	selectedCount := len(m.selected)
	b.WriteString(ErrorStyle.Render(fmt.Sprintf("Delete %d selected files?", selectedCount)))
	b.WriteString("\n\n")

	// Show first few files to be deleted
	count := 0
	for i, selected := range m.selected {
		if selected && count < 5 && i < len(m.results) {
			b.WriteString(PathStyle.Render(fmt.Sprintf("• %s (%s)",
				m.results[i].Path,
				utils.HumanSize(m.results[i].Size))))
			b.WriteString("\n")
			count++
		}
	}

	if selectedCount > 5 {
		b.WriteString(PathStyle.Render(fmt.Sprintf("... and %d more files", selectedCount-5)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(SuccessStyle.Render("y") + " to confirm, " + ErrorStyle.Render("n/esc") + " to cancel")
	return b.String()
}

func (m Model) helpView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("📖 Help"))
	b.WriteString("\n\n")
	b.WriteString(m.help.View(m.keys))
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("Press ? again to return"))
	return b.String()
}

func (m *Model) updateTable() {
	rows := make([]table.Row, len(m.results))
	for i, item := range m.results {
		selected := "[ ]"
		if m.selected[i] {
			selected = SelectedStyle.Render("[✓]")
		}
		rows[i] = table.Row{
			selected,
			SizeStyle.Render(utils.HumanSize(item.Size)),
			PathStyle.Render(item.Path),
		}
	}
	m.table.SetRows(rows)
}

func (m Model) hasSelected() bool {
	return len(m.selected) > 0
}

func (m Model) startScan() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		s := scanner.New(m.config)
		results, stats := s.Scan()
		return scanCompleteMsg{results: results, stats: stats}
	})
}

func (m Model) removeSelected() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		var removed, errors int

		for i, selected := range m.selected {
			if selected && i < len(m.results) {
				if err := removeFile(m.results[i].Path); err != nil {
					errors++
				} else {
					removed++
				}
			}
		}

		message := fmt.Sprintf("Removed %d files", removed)
		if errors > 0 {
			message += fmt.Sprintf(" (%d errors)", errors)
		}

		return removeCompleteMsg{
			removed: removed,
			errors:  errors,
			message: message,
		}
	})
}