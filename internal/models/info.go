package models

// Info структура для ответа API
type Info struct {
	Coins       int            `json:"coins"`
	Inventory   []PurchaseInfo `json:"inventory"`
	CoinHistory CoinHistory    `json:"coinHistory"`
}

type PurchaseInfo struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []TransferInfo `json:"received"`
	Sent     []TransferInfo `json:"sent"`
}

type TransferInfo struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}
