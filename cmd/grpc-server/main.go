package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/finebiscuit/server/services/balances"
	"github.com/finebiscuit/server/storage/inmem"
)

const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := inmem.New()
	balancesSvc := balances.NewService(db.BalancesTxFn())

	mux := http.NewServeMux()
	mux.Handle(balances.NewHandler(balancesSvc))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		panic(err)
	}
}
