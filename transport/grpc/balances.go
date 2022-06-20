package grpc

import (
	"context"
	"encoding/base64"

	pb "github.com/finebiscuit/proto/biscuit/accounting/v1"

	"github.com/finebiscuit/server/model/payload"
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

func (s *balancesServer) ListBalances(ctx context.Context, _ *pb.ListBalancesRequest) (*pb.ListBalancesResponse, error) {
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
	bp, err := protoToPayload(req.GetBalance().GetEncData())
	if err != nil {
		return nil, err
	}
	b, err := balance.New(req.Balance.TypeId, req.Balance.CurrencyId, bp)
	if err != nil {
		return nil, err
	}

	ep, err := protoToPayload(req.GetBalance().GetEncData())
	if err != nil {
		return nil, err
	}
	e, err := balance.NewEntry(ep)
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
		EncData:    payloadToProto(b.Payload),
		CurrentEntry: &pb.Entry{
			Id:      b.Entry.YMD.String(),
			EncData: payloadToProto(b.Entry.Payload),
		},
	}
}

func payloadToProto(p payload.Payload) *pb.EncryptedData {
	value := base64.StdEncoding.EncodeToString(p.Blob)
	return &pb.EncryptedData{
		VersionHash: p.Version,
		Algo:        uint32(p.Scheme),
		Base64Value: value,
	}
}

func protoToPayload(encData *pb.EncryptedData) (payload.Payload, error) {
	s, err := payload.NewScheme(int(encData.Algo))
	if err != nil {
		return payload.Payload{}, err
	}
	blob, err := base64.StdEncoding.DecodeString(encData.Base64Value)
	if err != nil {
		return payload.Payload{}, err
	}
	return payload.New(s, encData.VersionHash, blob)
}
