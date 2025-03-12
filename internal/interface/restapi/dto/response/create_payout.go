package response

type Sale struct {
	ID               string  `json:"id"`
	ItemID           string  `json:"itemId"`
	Price            float32 `json:"price"`
	OriginalCurrency string  `json:"originalCurrency"`
}

type Payout struct {
	ID     string  `json:"id"`
	Amount float32 `json:"amount"`
}

type BatchPayout struct {
	ID              string   `json:"id"`
	SellerReference string   `json:"sellerReference"`
	Payouts         []Payout `json:"payouts"`
	Sales           []Sale   `json:"sales"`
}

type CreatePayoutResponse struct {
	Currency     string        `json:"currency"`
	BatchPayouts []BatchPayout `json:"batchPayouts"`
}
