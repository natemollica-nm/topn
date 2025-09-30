package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginBottom(1).
		Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#10B981"))

	FileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F3F4F6"))

	SizeStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F59E0B"))

	PathStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))

	SelectedStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#374151")).
		Foreground(lipgloss.Color("#F9FAFB")).
		Bold(true).
		Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true)

	InfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6"))

	ProgressStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B5CF6"))

	WarningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Bold(true)

	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#374151")).
		Padding(1, 2)

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		MarginTop(1)
)