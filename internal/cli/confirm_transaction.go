package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/gustavo-bordin/thunes/internal/repository"
	"github.com/gustavo-bordin/thunes/internal/thunes"
	"go.mongodb.org/mongo-driver/bson"
)

type createTransactionRequest struct {
	isLoading   bool
	err         error
	transaction thunes.Transaction
}

type confirmTransactionRequest struct {
	isLoading   bool
	err         error
	transaction thunes.Transaction
}

type getTransactionStatesRequest struct {
	isLoading bool
	err       error
	states    []thunes.TransactionState
}

type SummaryScreen struct {
	thunesClient thunes.ThunesClient
	ngrokUrl     string

	confirmTransactionReq   confirmTransactionRequest
	createTransactionReq    createTransactionRequest
	getTransactionStatesReq getTransactionStatesRequest
	transactionRepo         repository.Repository
	externalId              string
	quotationId             int
	viewport                viewport.Model
	renderer                *glamour.TermRenderer
}

func NewSummaryScreen(
	tc thunes.ThunesClient,
	externalId string,
	quotationId int,
	tr repository.Repository,
	ngrokUrl string,
) SummaryScreen {
	return SummaryScreen{
		thunesClient:    tc,
		externalId:      externalId,
		quotationId:     quotationId,
		transactionRepo: tr,
		ngrokUrl:        ngrokUrl,
	}
}

func (s SummaryScreen) Init() tea.Cmd {
	s.createTransactionReq.isLoading = true
	return s.createTransaction
}

func (s SummaryScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyEnter:
			s.getTransactionStatesReq.isLoading = true
			return s, s.getTransactionStates
		}

	case createTransactionRequest:
		s.createTransactionReq.isLoading = false
		s.confirmTransactionReq.isLoading = true
		s.createTransactionReq.transaction = msg.transaction
		return s, s.confirmTransaction

	case confirmTransactionRequest:
		s.confirmTransactionReq.isLoading = false
		s.getTransactionStatesReq.isLoading = true
		return s, s.getTransactionStates

	case getTransactionStatesRequest:
		s.getTransactionStatesReq.isLoading = false
		s.getTransactionStatesReq.states = msg.states
		return s, nil
	}

	return s, nil
}

func (s SummaryScreen) View() string {
	if s.getTransactionStatesReq.err != nil {
		return errorStyle.Render(s.createTransactionReq.err.Error())
	}

	if s.getTransactionStatesReq.isLoading {
		loadingMsg := fmt.Sprintf("[%s] Loading states", s.externalId)
		return titleStyle.Render(loadingMsg)
	}

	if s.confirmTransactionReq.err != nil {
		return errorStyle.Render(s.createTransactionReq.err.Error())
	}

	if s.confirmTransactionReq.isLoading {
		loadingMsg := fmt.Sprintf("[%s] Confirming your transaction", s.externalId)
		return titleStyle.Render(loadingMsg)
	}

	if s.createTransactionReq.err != nil {
		return errorStyle.Render(s.createTransactionReq.err.Error())
	}

	if s.createTransactionReq.isLoading {
		loadingMsg := fmt.Sprintf("[%s] Creating your transaction", s.externalId)
		return titleStyle.Render(loadingMsg)
	}

	var msg strings.Builder

	titleMsg := fmt.Sprintf("[%s] Transaction done\n\n", s.externalId)
	msg.WriteString(titleStyle.Render(titleMsg))

	msg.WriteString(s.viewport.View())

	stateMsgFormat := "\n%s - %s\n"

	stateMsg := fmt.Sprintf(
		stateMsgFormat,
		s.createTransactionReq.transaction.CreationDate,
		s.createTransactionReq.transaction.StatusMessage,
	)
	msg.WriteString(stateMsg)

	for _, state := range s.getTransactionStatesReq.states {
		stateMsg := fmt.Sprintf(stateMsgFormat, state.CreationDate, state.StatusMessage)
		msg.WriteString(stateMsg)
	}
	msg.WriteString("\n\nPress ENTER to refresh states\n\n")
	msg.WriteString(getKeyMapMsg())

	return msg.String()
}

func (s SummaryScreen) createTransaction() tea.Msg {
	payload := thunes.CreateTransactionDto{
		CreditPartyIdentifier: thunes.TransactionIdentifier{
			MSISDN: "+263775892100",
		},
		Sender: thunes.TransactionEntityInfo{
			FirstName: "Gustavo",
			LastName:  "Bordin",
		},
		Beneficiary: thunes.TransactionEntityInfo{
			FirstName: "Christina",
			LastName:  "Lee",
		},
		ExternalID:          s.externalId,
		CallbackURL:         s.ngrokUrl,
		PurposeOfRemittance: "OTHER",
	}

	transaction, err := s.thunesClient.CreateTransaction(payload, s.quotationId)
	if err != nil {
		return createTransactionRequest{err: err}
	}

	return createTransactionRequest{transaction: *transaction, err: err}
}

func (s SummaryScreen) confirmTransaction() tea.Msg {
	transaction, err := s.thunesClient.ConfirmTransaction(s.externalId)
	if err != nil {
		return confirmTransactionRequest{err: err}
	}

	return confirmTransactionRequest{transaction: *transaction, err: err}
}

func (s SummaryScreen) getTransactionStates() tea.Msg {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	filter := bson.M{"externalid": s.externalId}
	res, err := s.transactionRepo.Find(ctx, filter)
	if err != nil {
		return getTransactionStatesRequest{err: err}
	}

	return getTransactionStatesRequest{states: res, err: err}
}
