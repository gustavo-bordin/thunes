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
	balancesLoadingMsg = "Loading balances"
	balancesErrMsg     = "unexpected error happened while loading balances: %s"
)

type balancesRequest struct {
	isLoading bool
	balances  []thunes.Balance
	err       error
}

type BalancesScreen struct {
	thunesClient thunes.ThunesClient
	ngrokUrl     string

	spinner spinner.Model
	cursor  int

	selectedPayer   thunes.Payer
	selectedBalance thunes.Balance
	balancesRequest balancesRequest
	transactionRepo repository.Repository
}

func NewBalancesScreen(
	tc thunes.ThunesClient,
	selectedPayer thunes.Payer,
	tr repository.Repository,
	ngrokUrl string,
) BalancesScreen {
	s := spinner.New()
	s.Spinner = spinner.Dot

	balancesRequest := balancesRequest{
		isLoading: true,
	}

	return BalancesScreen{
		spinner:         s,
		thunesClient:    tc,
		balancesRequest: balancesRequest,
		selectedPayer:   selectedPayer,
		transactionRepo: tr,
		ngrokUrl:        ngrokUrl,
	}
}

func (s BalancesScreen) getBalances() tea.Msg {
	balances, err := s.thunesClient.GetBalances()
	return balancesRequest{balances: balances, err: err}
}

func (s BalancesScreen) Init() tea.Cmd {
	return tea.Batch(s.spinner.Tick, s.getBalances)
}

func (s BalancesScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if s.cursor < len(s.balancesRequest.balances)-1 {
				s.cursor++
			}
			return s, nil

		case tea.KeyLeft:
			previousScreen := NewPayersScreen(s.thunesClient, s.transactionRepo, s.ngrokUrl)
			return NewRootScreen(s.thunesClient, s.transactionRepo, s.ngrokUrl).SwitchScreen(&previousScreen)

		case tea.KeyEnter:
			s.selectedBalance = s.balancesRequest.balances[s.cursor]
			nextScreen := NewAmountScreen(
				s.thunesClient,
				s.selectedPayer,
				s.selectedBalance,
				s.transactionRepo,
				s.ngrokUrl,
			)
			return NewRootScreen(s.thunesClient, s.transactionRepo, s.ngrokUrl).SwitchScreen(&nextScreen)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd

	case balancesRequest:
		s.balancesRequest.balances = msg.balances
		s.balancesRequest.err = msg.err
		s.balancesRequest.isLoading = false
		return s, nil
	}

	return s, nil
}

func (s BalancesScreen) View() string {
	if s.balancesRequest.err != nil {
		err := fmt.Sprintf(balancesErrMsg, s.balancesRequest.err)
		return err
	}

	if s.balancesRequest.isLoading {
		loadingSpinner := s.spinner.View() + " "
		return loadingSpinner + balancesLoadingMsg
	}

	var sb strings.Builder

	title := titleStyle.Render("Please select one of the balances below")
	sb.WriteString(title)
	sb.WriteString("\n\n\n")

	for index, balance := range s.balancesRequest.balances {
		if index == s.cursor {
			sb.WriteString("[💰] ")
		} else {
			sb.WriteString("[ ] ")
		}

		balanceInfo := fmt.Sprintf(
			"%.2f %s in %s",
			balance.Balance,
			balance.Currency,
			balance.Name,
		)

		balanceMsg := balanceInfoStyle.Render(balanceInfo)
		sb.WriteString(balanceMsg)
	}

	sb.WriteString(getKeyMapMsg())

	return sb.String()
}
