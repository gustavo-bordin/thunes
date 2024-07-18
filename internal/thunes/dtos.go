package thunes

type CreateQuotationDto struct {
	ExternalId      string           `json:"external_id"`
	PayerId         int32            `json:"payer_id"`
	Mode            string           `json:"mode"`
	TransactionType string           `json:"transaction_type"`
	Source          TransactionMoney `json:"source"`
	Destination     TransactionMoney `json:"destination"`
}

type CreateTransactionDto struct {
	CreditPartyIdentifier TransactionIdentifier `json:"credit_party_identifier"`
	Sender                TransactionEntityInfo `json:"sender"`
	Beneficiary           TransactionEntityInfo `json:"beneficiary"`
	ExternalID            string                `json:"external_id"`
	CallbackURL           string                `json:"callback_url"`
	PurposeOfRemittance   string                `json:"purpose_of_remittance"`
}
