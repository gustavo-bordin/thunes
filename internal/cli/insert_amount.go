package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gustavo-bordin/thunes/internal/repository"
	"github.com/gustavo-bordin/thunes/internal/thunes"
)

type AmountScreen struct {
	thunesClient thunes.ThunesClient
	ngrokUrl     string

	selectedPayer   thunes.Payer
	selectedBalance thunes.Balance
	amountInput     textinput.Model
	transactionRepo repository.Repository

	err error
}

func NewAmountScreen(
	tc thunes.ThunesClient,
	selectedPayer thunes.Payer,
	selectedBalance thunes.Balance,
	tr repository.Repository,
	ngrokUrl string,
) AmountScreen {
	amountInput := textinput.New()
	amountInput.Placeholder = "Insert amount here"
	amountInput.Focus()

	return AmountScreen{
		thunesClient:    tc,
		selectedPayer:   selectedPayer,
		selectedBalance: selectedBalance,
		amountInput:     amountInput,
		transactionRepo: tr,
		ngrokUrl:        ngrokUrl,
	}
}

func (s AmountScreen) parseAmount() (float64, error) {
	value := s.amountInput.Value()

	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return amount, err
	}

	return amount, nil
}

func (s AmountScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s AmountScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyLeft:
			previousScreen := NewBalancesScreen(s.thunesClient, s.selectedPayer, s.transactionRepo, s.ngrokUrl)
			return NewRootScreen(s.thunesClient, s.transactionRepo, s.ngrokUrl).SwitchScreen(&previousScreen)
		case tea.KeyEnter:
			amount, err := s.parseAmount()
			if err != nil {
				s.err = err
				return s, nil
			}

			nextScreen := NewConfirmScreen(
				s.thunesClient,
				s.selectedPayer,
				s.selectedBalance,
				s.transactionRepo,
				amount,
				s.ngrokUrl,
			)
			return NewRootScreen(s.thunesClient, s.transactionRepo, s.ngrokUrl).SwitchScreen(&nextScreen)
		}
	case error:
		s.err = msg
		return s, nil
	}

	var cmd tea.Cmd
	s.amountInput, cmd = s.amountInput.Update(msg)
	return s, cmd
}

func (s AmountScreen) View() string {
	var msg strings.Builder

	payerName := payerNameStyle.Render(s.selectedPayer.Name)
	title := fmt.Sprintf("Sending money to %s", payerName)
	titleMsg := titleStyle.Render(title)
	msg.WriteString(titleMsg)
	msg.WriteString("\n\n")

	balanceMsg := fmt.Sprintf(
		"⚠️ You have %.2f %s available in %s",
		s.selectedBalance.Balance,
		s.selectedBalance.Currency,
		s.selectedBalance.Name,
	)

	boundariesMsg := fmt.Sprintf(
		"⚠️ Min of %s and Max of %s %s",
		s.selectedPayer.TransactionTypes.C2C.MinimumTransactionAmount,
		s.selectedPayer.TransactionTypes.C2C.MaximumTransactionAmount,
		s.selectedPayer.Currency,
	)

	msg.WriteString(warningStyle.Render(boundariesMsg))
	msg.WriteString("\n")
	msg.WriteString(warningStyle.Render(balanceMsg))

	if s.err != nil {
		msg.WriteString("\n")
		errMsg := errorStyle.Render("❌ Value is not valid amount")
		msg.WriteString(errMsg)

	}

	msg.WriteString("\n\n\n")

	inputMsg := fmt.Sprintf(
		"Insert amount in %s %s",
		s.selectedBalance.Currency,
		s.amountInput.View(),
	)

	msg.WriteString(inputMsg)
	msg.WriteString(getKeyMapMsg())

	return msg.String()
}
