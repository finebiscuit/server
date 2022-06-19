package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	balancesPb "github.com/finebiscuit/proto/biscuit/accounting/v1"
	"github.com/finebiscuit/server/services/balances"
	"github.com/finebiscuit/server/storage/inmem"
	transport "github.com/finebiscuit/server/transport/grpc"
)

const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := inmem.New()
	balancesSvc := balances.NewService(db.BalancesTxFn())
	balancesServer := transport.NewBalancesServer(balancesSvc)

	grpcSrv := grpc.NewServer(
	// TODO: add production-grade interceptors
	)
	balancesPb.RegisterAccountingServer(grpcSrv, balancesServer)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	if err := grpcSrv.Serve(l); err != nil {
		panic(err)
	}
}
