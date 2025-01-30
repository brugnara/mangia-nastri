package logger

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// New initializes and returns a customized logger instance.
// It sets up specific styles for log levels and keys using the lipgloss
// package to enhance the visual appearance of log messages.
//
// Returns
//   A logger instance with the specified styles applied.

func New(context string) log.Logger {
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("220")).
		Foreground(lipgloss.Color("#FFFFFF"))

	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	logger := log.New(os.Stderr)
	logger.SetStyles(styles)
	logger.SetPrefix(context)

	return *logger
}
