package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/gustavo-bordin/thunes/internal/repository"
	"github.com/gustavo-bordin/thunes/internal/thunes"
)

type quotationRequest struct {
	isLoading bool
	err       error
	quotation thunes.Quotation
}

type ConfirmScreen struct {
	thunesClient thunes.ThunesClient
	ngrokUrl     string

	spinner spinner.Model

	selectedPayer    thunes.Payer
	selectedBalance  thunes.Balance
	quotationRequest quotationRequest
	transactionRepo  repository.Repository
	externalId       string
	amount           float64
}

func NewConfirmScreen(
	tc thunes.ThunesClient,
	selectedPayer thunes.Payer,
	selectedBalance thunes.Balance,
	tr repository.Repository,
	amount float64,
	ngrokUrl string,
) ConfirmScreen {

	s := spinner.New()
	s.Spinner = spinner.Dot

	quotationRequest := quotationRequest{isLoading: true}

	externalId := uuid.New().String()

	return ConfirmScreen{
		spinner:          s,
		thunesClient:     tc,
		selectedPayer:    selectedPayer,
		selectedBalance:  selectedBalance,
		quotationRequest: quotationRequest,
		externalId:       externalId,
		amount:           amount,
		transactionRepo:  tr,
		ngrokUrl:         ngrokUrl,
	}
}

func (s ConfirmScreen) createQuotation() tea.Msg {
	quotation := thunes.CreateQuotationDto{
		ExternalId:      s.externalId,
		PayerId:         s.selectedPayer.ID,
		Mode:            "SOURCE_AMOUNT",
		TransactionType: "C2C",
		Source: thunes.TransactionMoney{
			Amount:         &s.amount,
			Currency:       s.selectedBalance.Currency,
			CountryIsoCode: "BRA",
		},
		Destination: thunes.TransactionMoney{
			Currency:       s.selectedPayer.Currency,
			CountryIsoCode: s.selectedPayer.CountryIsoCode,
		},
	}

	res, err := s.thunesClient.CreateQuotation(quotation)
	return quotationRequest{quotation: *res, err: err}
}

func (s ConfirmScreen) Init() tea.Cmd {
	return tea.Batch(s.spinner.Tick, s.createQuotation)
}

func (s ConfirmScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit

		case tea.KeyEnter:
			quotationId := s.quotationRequest.quotation.ID
			nextScreen := NewSummaryScreen(
				s.thunesClient,
				s.externalId,
				quotationId,
				s.transactionRepo,
				s.ngrokUrl,
			)
			return NewRootScreen(s.thunesClient, s.transactionRepo, s.ngrokUrl).SwitchScreen(&nextScreen)

		case tea.KeyLeft:
			previousScreen := NewAmountScreen(
				s.thunesClient,
				s.selectedPayer,
				s.selectedBalance,
				s.transactionRepo,
				s.ngrokUrl,
			)
			return NewRootScreen(s.thunesClient, s.transactionRepo, s.ngrokUrl).SwitchScreen(&previousScreen)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd

	case quotationRequest:
		s.quotationRequest.err = msg.err
		s.quotationRequest.isLoading = false
		s.quotationRequest.quotation = msg.quotation
		return s, nil
	}

	return s, nil
}

func (s ConfirmScreen) View() string {
	var msg strings.Builder

	payerName := payerNameStyle.Render(s.selectedPayer.Name)
	title := fmt.Sprintf(
		"You are about to send %.2f %s to %s",
		s.amount,
		s.selectedBalance.Currency,
		payerName,
	)
	titleMsg := titleStyle.Render(title)
	msg.WriteString(titleMsg)
	msg.WriteString("\n\n")

	if s.quotationRequest.isLoading {
		msg.WriteString(s.spinner.View() + " Loading quotation")
		msg.WriteString("\n\n")
		return msg.String()
	}

	hasQuotationErrors := len(s.quotationRequest.quotation.Errors) != 0

	if hasQuotationErrors {
		for _, err := range s.quotationRequest.quotation.Errors {
			errorMsg := fmt.Sprintf("❌ [%s] %s", err.Code, err.Message)
			msg.WriteString(errorMsg)
			msg.WriteString("\n")
		}

		msg.WriteString(getKeyMapMsg())
		return msg.String()
	}

	feeMsg := fmt.Sprintf(
		"⚠️ You will be charged %.2f %s for this transaction",
		*s.quotationRequest.quotation.Fee.Amount,
		s.selectedBalance.Currency,
	)
	msg.WriteString(warningStyle.Render(feeMsg))
	msg.WriteString("\n")

	conversionMsg := fmt.Sprintf(
		"⚠️ %s will receive %.2f %s",
		s.selectedPayer.Name,
		*s.quotationRequest.quotation.Destination.Amount,
		s.selectedPayer.Currency,
	)

	msg.WriteString(warningStyle.Render(conversionMsg))
	msg.WriteString("\n\n")

	msg.WriteString("Press ENTER to confirm your transaction")

	msg.WriteString(getKeyMapMsg())

	return msg.String()
}
