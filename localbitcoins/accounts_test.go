package localbitcoins

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccount_marshall(t *testing.T) {
	testJSONMarshal(t, &Account{}, "{}")

	a := &Account{
		Username:                  String("foo"),
		TradingPartnersCount:      Int(5),
		FeedbacksUnconfirmedCount: Int(2),
		TradeVolumeText:           String("Less than 25 BTC"),
		HasCommonTrades:           Bool(false),
		ConfirmedTradeCountText:   String("0"),
		BlockedCount:              Int(0),
		FeedbackCount:             Int(0),
		Url:                       String("https://localbitcoins.com/p/foo/"),
		TrustedCount:              Int(2),
	}
	want := `{
    "username": "foo",
    "trading_partners_count": 5,
    "feedbacks_unconfirmed_count": 2,
    "trade_volume_text": "Less than 25 BTC",
    "has_common_trades": false,
    "confirmed_trade_count_text": "0",
    "blocked_count": 0,
    "feedback_count": 0,
    "url": "https://localbitcoins.com/p/foo/",
    "trusted_count": 2
  }`
	testJSONMarshal(t, a, want)
}

func TestAccountsService_Get_specifiedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/account_info/foo/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data":{"username":"foo"}}`)
	})

	acc, _, err := client.Accounts.Get("foo")
	if err != nil {
		t.Errorf("Accounts.Get returned error: %v", err)
	}

	want := &Account{Username: String("foo")}
	if !reflect.DeepEqual(acc, want) {
		t.Errorf("Accounts.Get returned %+v, want %+v", acc, want)
	}
}

func TestAccountsService_Get_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/myself/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data":{"username": "foo"}}`)
	})

	acc, _, err := client.Accounts.Get("")
	if err != nil {
		t.Errorf("Accounts.Get returned error: %v", err)
	}

	want := &Account{Username: String("foo")}
	if !reflect.DeepEqual(acc, want) {
		t.Errorf("Accounts.Get returned %+v, want %+v", acc, want)
	}
}
