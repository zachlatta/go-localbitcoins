package localbitcoins

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccountsService_Get_specifiedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/account_info/foo", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("/api/myself", func(w http.ResponseWriter, r *http.Request) {
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
