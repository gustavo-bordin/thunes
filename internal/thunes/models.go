package thunes

import "time"

type ThunesError struct {
	Code    string
	Message string
}

type Balance struct {
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type Payer struct {
	ID               int32                 `json:"id"`
	Name             string                `json:"name"`
	Currency         string                `json:"currency"`
	CountryIsoCode   string                `json:"country_iso_code"`
	TransactionTypes PayerTransactionTypes `json:"transaction_types"`
}

type PayerTransactionTypesC2C struct {
	MinimumTransactionAmount string `json:"minimum_transaction_amount"`
	MaximumTransactionAmount string `json:"maximum_transaction_amount"`
}

type PayerTransactionTypes struct {
	C2C PayerTransactionTypesC2C `json:"C2C"`
}

type TransactionMoney struct {
	Amount         *float64 `json:"amount"`
	Currency       string   `json:"currency"`
	CountryIsoCode string   `json:"country_iso_code"`
}

type Quotation struct {
	ID          int
	Destination TransactionMoney
	Fee         TransactionMoney
	SentAmount  TransactionMoney
	Errors      []ThunesError
}

type Transaction struct {
	Beneficiary           TransactionEntityInfo `json:"beneficiary"`
	CallbackURL           string                `json:"callback_url"`
	CreationDate          time.Time             `json:"creation_date"`
	CreditPartyIdentifier TransactionIdentifier `json:"credit_party_identifier"`
	Destination           TransactionMoney      `json:"destination"`
	ExpirationDate        time.Time             `json:"expiration_date"`
	ExternalID            string                `json:"external_id"`
	Fee                   TransactionMoney      `json:"fee"`
	ID                    int                   `json:"id"`
	Payer                 Payer                 `json:"payer"`
	PurposeOfRemittance   string                `json:"purpose_of_remittance"`
	Sender                TransactionEntityInfo `json:"sender"`
	SentAmount            TransactionMoney      `json:"sent_amount"`
	Source                TransactionMoney      `json:"source"`
	Status                string                `json:"status"`
	StatusClass           string                `json:"status_class"`
	StatusClassMessage    string                `json:"status_class_message"`
	StatusMessage         string                `json:"status_message"`
	TransactionType       string                `json:"transaction_type"`
	WholesaleFxRate       float64               `json:"wholesale_fx_rate"`
	Errors                []ThunesError
}

type TransactionState struct {
	Beneficiary           TransactionEntityInfo `json:"beneficiary"`
	CallbackURL           string                `json:"callback_url"`
	CreationDate          time.Time             `json:"creation_date"`
	CreditPartyIdentifier TransactionIdentifier `json:"credit_party_identifier"`
	Destination           TransactionMoney      `json:"destination"`
	ExpirationDate        time.Time             `json:"expiration_date"`
	ExternalID            string                `json:"external_id"`
	Fee                   TransactionMoney      `json:"fee"`
	ID                    int                   `json:"id"`
	Payer                 Payer                 `json:"payer"`
	PurposeOfRemittance   string                `json:"purpose_of_remittance"`
	Sender                TransactionEntityInfo `json:"sender"`
	SentAmount            TransactionMoney      `json:"sent_amount"`
	Source                TransactionMoney      `json:"source"`
	Status                int                   `json:"status"`
	StatusClass           int                   `json:"status_class"`
	StatusClassMessage    string                `json:"status_class_message"`
	StatusMessage         string                `json:"status_message"`
	TransactionType       string                `json:"transaction_type"`
	WholesaleFxRate       float64               `json:"wholesale_fx_rate"`
	Errors                []ThunesError
}

type TransactionEntityInfo struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type TransactionIdentifier struct {
	MSISDN string `json:"msisdn"`
}
