package domain

type Currency string
type Market string

type Order struct {
	Market    string
	Ask       float64
	Bid       float64
	Timestamp int64
}
