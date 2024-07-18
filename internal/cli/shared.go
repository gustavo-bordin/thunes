package cli

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

const (
	titleColor     = lipgloss.Color("#008e00")
	warningColor   = lipgloss.Color("#EED202")
	highlightColor = lipgloss.Color("#EE4B2B")

	errorColor = lipgloss.Color("#FF9494")
)

var (
	titleStyle       = lipgloss.NewStyle().Foreground(titleColor)
	payerNameStyle   = lipgloss.NewStyle().Foreground(highlightColor)
	balanceInfoStyle = lipgloss.NewStyle().Foreground(highlightColor)
	warningStyle     = lipgloss.NewStyle().Foreground(warningColor)
	errorStyle       = lipgloss.NewStyle().Foreground(errorColor)
)

type keyMapHelper struct {
	next, prev, up, down, confirm, quit key.Binding
}

func getKeyMapMsg() string {
	keyMapHelper := keyMapHelper{
		prev: key.NewBinding(
			key.WithKeys("←"),
			key.WithHelp("←", "prev"),
		),
		next: key.NewBinding(
			key.WithKeys("→"),
			key.WithHelp("→", "next"),
		),
		up: key.NewBinding(
			key.WithKeys("↑"),
			key.WithHelp("↑", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("↓"),
			key.WithHelp("↓", "down"),
		),
		confirm: key.NewBinding(
			key.WithKeys("↵"),
			key.WithHelp("↵", "confirm"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}

	helpKeys := help.New().ShortHelpView([]key.Binding{
		keyMapHelper.next,
		keyMapHelper.prev,
		keyMapHelper.up,
		keyMapHelper.down,
		keyMapHelper.confirm,
		keyMapHelper.quit,
	})

	kepMapMsg := "\n\n\n\n" + helpKeys
	return kepMapMsg
}
