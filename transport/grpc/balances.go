package grpc

import (
	"context"
	"encoding/base64"

	pb "github.com/finebiscuit/proto/biscuit/balances/v1"

	"github.com/finebiscuit/server/model/payload"
	"github.com/finebiscuit/server/services/balances"
	"github.com/finebiscuit/server/services/balances/balance"
)

type balancesServer struct {
	pb.UnimplementedBalancesServer
	Balances balances.Service
}

func NewBalancesServer(balancesService balances.Service) pb.BalancesServer {
	return &balancesServer{
		Balances: balancesService,
	}
}

func (s *balancesServer) ListBalances(ctx context.Context, _ *pb.ListRequest) (*pb.ListResponse, error) {
	bals, err := s.Balances.ListBalances(ctx, balance.Filter{})
	if err != nil {
		return nil, err
	}

	res := &pb.ListResponse{
		Balances: make([]*pb.Balance, len(bals)),
	}

	for _, b := range bals {
		res.Balances = append(res.Balances, balanceToProto(b))
	}
	return res, nil
}

func (s *balancesServer) GetBalance(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	id, err := balance.ParseID(req.GetBalanceId())
	if err != nil {
		return nil, err
	}

	b, err := s.Balances.GetBalance(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &pb.GetResponse{Balance: balanceToProto(b)}
	return res, nil
}

func (s *balancesServer) CreateBalance(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	bp, err := protoToPayload(req.GetBalancePayload())
	if err != nil {
		return nil, err
	}
	b, err := balance.New(req.GetTypeId(), req.GetCurrencyId(), bp)
	if err != nil {
		return nil, err
	}

	ep, err := protoToPayload(req.GetEntryPayload())
	if err != nil {
		return nil, err
	}
	e, err := balance.NewEntryWithString(req.GetEntryYmd(), ep)
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

	res := &pb.CreateResponse{Balance: balanceToProto(bwe)}
	return res, nil
}

func balanceToProto(b *balance.WithEntry) *pb.Balance {
	return &pb.Balance{
		Id:         b.ID.String(),
		TypeId:     b.TypeID,
		CurrencyId: b.CurrencyID,
		Payload:    payloadToProto(b.Payload),
		CurrentEntry: &pb.Entry{
			Ymd:     b.Entry.YMD.String(),
			Payload: payloadToProto(b.Entry.Payload),
		},
	}
}

func payloadToProto(p payload.Payload) *pb.Payload {
	value := base64.StdEncoding.EncodeToString(p.Blob)
	return &pb.Payload{
		Version:     p.Version,
		Scheme:      uint32(p.Scheme),
		Base64Value: value,
	}
}

func protoToPayload(p *pb.Payload) (payload.Payload, error) {
	s, err := payload.NewScheme(int(p.Scheme))
	if err != nil {
		return payload.Payload{}, err
	}
	blob, err := base64.StdEncoding.DecodeString(p.Base64Value)
	if err != nil {
		return payload.Payload{}, err
	}
	return payload.New(s, p.Version, blob)
}
