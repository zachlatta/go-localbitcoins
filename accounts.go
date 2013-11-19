package localbitcoins

import "fmt"

// AccountsService provides access to all account related functions in the
// LocalBitcoins API.
type AccountsService struct {
	client *Client
}

// Account represents a LocalBitcoins account.
type Account struct {
	Username                  *string `json:"username,omitempty"`
	TradingPartnersCount      *int    `json:"trading_partners_count,omitempty"`
	FeedbacksUnconfirmedCount *int    `json:"feedbacks_unconfirmed_count,omitempty"`
	TradeVolumeText           *string `json:"trade_volume_text,omitempty"`
	HasCommonTrades           *bool   `json:"has_common_trades,omitempty"`
	ConfirmedTradeCountText   *string `json:"confirmed_trade_count_text,omitempty"`
	BlockedCount              *int    `json:"blocked_count,omitempty"`
	FeedbackCount             *int    `json:"feedback_count,omitempty"`
	Url                       *string `json:"url,omitempty"`
	TrustedCount              *int    `json:"trusted_count,omitempty"`
}

func (a Account) String() string {
	return Stringify(a)
}

// Get fetches an account. Passing an empty string will fetch the authenticated
// account.
func (s *AccountsService) Get(account string) (*Account, *Response, error) {
	var a string
	if account != "" {
		a = fmt.Sprintf("api/account_info/%v", account)
	} else {
		a = "api/myself"
	}
	req, err := s.client.NewRequest("GET", a, nil)
	if err != nil {
		return nil, nil, err
	}

	acc := new(Account)
	aResp := ResponseData{Data: acc}
	resp, err := s.client.Do(req, &aResp)
	if err != nil {
		return nil, resp, err
	}

	return acc, resp, err
}
