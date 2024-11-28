package src

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/lipgloss"
)

func SetupLogger() log.Logger {
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("MangiaNastri").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("201")).
		Foreground(lipgloss.Color("#FFFFFF"))

	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	logger := log.New(os.Stderr)
	logger.SetStyles(styles)

	return *logger
}
