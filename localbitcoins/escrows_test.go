package localbitcoins

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestEscrow_Release(t *testing.T) {
	setup()
	defer teardown()

	escrow := &Escrow{
		releaseUrl: String("https://localbitcoins.com/api/escrow_release/1"),
	}

	mux.HandleFunc("/api/escrow_release/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `Success!`)
	})

	_, err := escrow.Release()
	if err != nil {
		t.Errorf("Error releasing escrow: %v", err)
	}
}

func TestEscrowsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/escrows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
      "data":{
        "escrow_list":[
        {
          "data":{"buyer_username":"foo"},
          "actions":{"release_url":"bar"}
        }
        ]
      }
    }`)
	})

	escrow, _, err := client.Escrows.List()
	if err != nil {
		t.Errorf("Escrows.List returned error: %v", err)
	}

	want := []*Escrow{
		&Escrow{BuyerUsername: String("foo"), releaseUrl: String("bar")},
	}
	if !reflect.DeepEqual(escrow, want) {
		t.Errorf("Escrows.List returned %+v, want %+v", escrow, want)
	}
}
