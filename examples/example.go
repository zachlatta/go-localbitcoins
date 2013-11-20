package main

import (
  "fmt"
  "github.com/zachlatta/go-localbitcoins/localbitcoins"
)

func main() {
  client := localbitcoins.NewClient(nil)
  acc, _, err := client.Accounts.Get("zrl")
  if err != nil {
    panic(err)
  }
  fmt.Println(acc)
}
