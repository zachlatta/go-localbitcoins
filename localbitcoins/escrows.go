package localbitcoins

import "time"

// EscrowsService handles all escrow-related communications with the
// LocalBitcoins API.
type EscrowsService struct {
	client *Client
}

// Escrow represents an escrow as returned by the LocalBitcoins API.
type Escrow struct {
	CreatedAt             *time.Time `json:"created_at,omitempty"`
	BuyerUsername         *string    `json:"buyer_username,omitempty"`
	ReferenceCode         *string    `json:"reference_code,omitempty"`
	Currency              *string    `json:"currency,omitempty"`
	Amount                *float64   `json:"amount,string,omitempty"`
	AmountBTC             *float64   `json:"amount_btc,string,omitempty"`
	ExchangeRateUpdatedAt *time.Time `json:"exchange_rate_updated_at,omitempty"`

	releaseUrl *string
}

func (e Escrow) String() string {
	return Stringify(e)
}

// Escrow list middleman used strictly for unmarshaling the API response.
type escrowListMiddleman struct {
	Escrows []*escrowMiddleman `json:"escrow_list,omitempty"`
}

// Middleman used strictly for unmarshaling individual escrows.
type escrowMiddleman struct {
	Escrow              *Escrow              `json:"data,omitempty"`
	ReleaseUrlMiddleman *releaseUrlMiddleman `json:"actions,omitempty"`
}

// Release URL middleman used solely for unmarshaling the release URL for an
// escrow.
type releaseUrlMiddleman struct {
	ReleaseUrl *string `json:"release_url,omitempty"`
}

func (s *EscrowsService) List() ([]*Escrow, *Response, error) {
	req, err := s.client.NewRequest("GET", "/api/escrows", nil)
	if err != nil {
		return nil, nil, err
	}

	middleman := new(escrowListMiddleman)
	respMiddleman := &ResponseData{Data: middleman}
	resp, err := s.client.Do(req, respMiddleman)
	if err != nil {
		return nil, resp, err
	}

	escrows := make([]*Escrow, len(middleman.Escrows))

	for i, e := range middleman.Escrows {
		escrows[i] = e.Escrow
		escrows[i].releaseUrl = e.ReleaseUrlMiddleman.ReleaseUrl
	}

	return escrows, resp, err
}
