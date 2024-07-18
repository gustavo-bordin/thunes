package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gustavo-bordin/thunes/internal/repository"
	"github.com/gustavo-bordin/thunes/internal/thunes"
)

var (
	availablePayersLoadingMsg = "Loading available payers"

	availablePayersErrMsg = "unexpected error happened while loading payers: %s"
)

type availablePayers struct {
	isLoading bool
	payers    []thunes.Payer
	err       error
}

type PayersScreen struct {
	thunesClient thunes.ThunesClient
	ngrokUrl     string

	spinner spinner.Model
	cursor  int

	selectedPayer   thunes.Payer
	availablePayers availablePayers
	transactionRepo repository.Repository
}

func NewPayersScreen(
	tc thunes.ThunesClient,
	tr repository.Repository,
	ngrokUrl string,
) PayersScreen {
	s := spinner.New()
	s.Spinner = spinner.Dot

	availablePayers := availablePayers{
		isLoading: true,
	}

	return PayersScreen{
		spinner:         s,
		thunesClient:    tc,
		availablePayers: availablePayers,
		transactionRepo: tr,
		ngrokUrl:        ngrokUrl,
	}
}

func (s PayersScreen) getAvailablePayers() tea.Msg {
	payers, err := s.thunesClient.GetAvailablePayers()
	return availablePayers{payers: payers, err: err}
}

func (s PayersScreen) Init() tea.Cmd {
	return tea.Batch(s.spinner.Tick, s.getAvailablePayers)
}

func (s PayersScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit

		case tea.KeyUp:
			if s.cursor > 0 {
				s.cursor--
			}
			return s, nil

		case tea.KeyDown:
			if s.cursor < len(s.availablePayers.payers)-1 {
				s.cursor++
			}
			return s, nil

		case tea.KeyEnter:
			s.selectedPayer = s.availablePayers.payers[s.cursor]
			nextScreen := NewBalancesScreen(s.thunesClient, s.selectedPayer, s.transactionRepo, s.ngrokUrl)
			return NewRootScreen(s.thunesClient, s.transactionRepo, s.ngrokUrl).SwitchScreen(&nextScreen)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd

	case availablePayers:
		s.availablePayers.payers = msg.payers
		s.availablePayers.err = msg.err
		s.availablePayers.isLoading = false
		return s, nil
	}

	return s, nil
}

func (s PayersScreen) View() string {
	if s.availablePayers.err != nil {
		err := fmt.Sprintf(availablePayersErrMsg, s.availablePayers.err)
		return err
	}

	if s.availablePayers.isLoading {
		loadingSpinner := s.spinner.View() + " "
		return loadingSpinner + availablePayersLoadingMsg
	}

	var sb strings.Builder

	title := titleStyle.Render("Please select one of the payers below")
	sb.WriteString(title)
	sb.WriteString("\n\n\n")

	for index, payer := range s.availablePayers.payers {
		if index == s.cursor {
			sb.WriteString("[💰] ")
		} else {
			sb.WriteString("[ ] ")
		}

		payerName := payerNameStyle.Render(payer.Name)
		sb.WriteString(payerName)
	}

	sb.WriteString(getKeyMapMsg())

	return sb.String()
}
