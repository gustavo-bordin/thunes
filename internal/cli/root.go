package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gustavo-bordin/thunes/internal/repository"
	"github.com/gustavo-bordin/thunes/internal/thunes"
)

type RootScreen struct {
	thunesClient    thunes.ThunesClient
	ngrokUrl        string
	currentScreen   tea.Model
	transactionRepo repository.Repository
}

func NewRootScreen(
	tc thunes.ThunesClient,
	r repository.Repository,
	ngrokUrl string,
) RootScreen {
	currentScreen := NewPayersScreen(tc, r, ngrokUrl)

	return RootScreen{
		currentScreen:   currentScreen,
		thunesClient:    tc,
		transactionRepo: r,
		ngrokUrl:        ngrokUrl,
	}
}

func (r RootScreen) Init() tea.Cmd {
	return r.currentScreen.Init()
}

func (r RootScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return r.currentScreen.Update(msg)
}

func (r RootScreen) View() string {
	return r.currentScreen.View()
}

func (r RootScreen) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	r.currentScreen = model
	return r.currentScreen, r.currentScreen.Init()
}
