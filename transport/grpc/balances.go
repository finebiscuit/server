package grpc

import (
	"context"

	pb "github.com/finebiscuit/proto/biscuit/accounting/v1"
	"github.com/finebiscuit/server/services/balances"
	"github.com/finebiscuit/server/services/balances/balance"
)

type balancesServer struct {
	pb.UnimplementedAccountingServer
	Balances balances.Service
}

func NewBalancesServer(balancesService balances.Service) pb.AccountingServer {
	return &balancesServer{
		Balances: balancesService,
	}
}

func (s *balancesServer) ListBalances(ctx context.Context, req *pb.ListBalancesRequest) (*pb.ListBalancesResponse, error) {
	bals, err := s.Balances.ListBalances(ctx, balance.Filter{})
	if err != nil {
		return nil, err
	}

	res := &pb.ListBalancesResponse{
		Balances: make([]*pb.Balance, len(bals)),
	}

	for _, b := range bals {
		res.Balances = append(res.Balances, balanceToProto(b))
	}
	return res, nil
}

func (s *balancesServer) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	id, err := balance.ParseID(req.GetBalanceId())
	if err != nil {
		return nil, err
	}

	b, err := s.Balances.GetBalance(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &pb.GetBalanceResponse{Balance: balanceToProto(b)}
	return res, nil
}

func (s *balancesServer) CreateBalance(ctx context.Context, req *pb.CreateBalanceRequest) (*pb.CreateBalanceResponse, error) {
	b, err := balance.New(req.Balance.TypeId, req.Balance.CurrencyId)
	if err != nil {
		return nil, err
	}

	e, err := balance.NewEntry()
	if err != nil {
		return nil, err
	}

	if err := s.Balances.CreateBalance(ctx, b, e); err != nil {
		return nil, err
	}

	bwe, err := s.Balances.GetBalance(ctx, b.ID)
	if err != nil {
		return nil, err
	}

	res := &pb.CreateBalanceResponse{Balance: balanceToProto(bwe)}
	return res, nil
}

func balanceToProto(b *balance.WithEntry) *pb.Balance {
	return &pb.Balance{
		Id:         b.ID.String(),
		TypeId:     b.TypeID,
		CurrencyId: b.CurrencyID,
		CurrentEntry: &pb.Entry{
			Id: b.Entry.YMD.String(),
		},
	}
}
