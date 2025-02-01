package logger

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Logger struct {
	log.Logger
}

// New initializes and returns a customized logger instance.
// It sets up specific styles for log levels and keys using the lipgloss
// package to enhance the visual appearance of log messages.
//
// Returns
//   A logger instance with the specified styles applied.

func New(context string) Logger {
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

	return Logger{*logger}
}

func (l Logger) CloneWithPrefix(prefix string) Logger {
	return New(l.GetPrefix() + ":" + prefix)
}
