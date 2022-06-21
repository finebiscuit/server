package balances

import (
	"context"

	"github.com/bufbuild/connect-go"
	balancesv1 "github.com/finebiscuit/proto/biscuit/balances/v1"
	"github.com/finebiscuit/proto/biscuit/balances/v1/balancesv1connect"

	"github.com/finebiscuit/server/services/balances/balance"
)

func NewClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) Service {
	return &client{
		Balances: balancesv1connect.NewBalancesClient(httpClient, baseURL, opts...),
	}
}

type client struct {
	Balances balancesv1connect.BalancesClient
}

var _ Service = &client{}

func (c *client) GetBalance(ctx context.Context, id balance.ID) (*balance.WithEntry, error) {
	req := connect.NewRequest(&balancesv1.GetRequest{BalanceId: id.String()})
	res, err := c.Balances.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	bwe, err := balance.NewWithEntryFromProto(res.Msg.GetBalance())
	if err != nil {
		return nil, err
	}
	return bwe, nil
}

func (c *client) ListBalances(ctx context.Context, filter balance.Filter) ([]*balance.WithEntry, error) {
	req := connect.NewRequest(&balancesv1.ListRequest{})
	res, err := c.Balances.List(ctx, req)
	if err != nil {
		return nil, err
	}

	bals := make([]*balance.WithEntry, 0, len(res.Msg.GetBalances()))
	for _, proto := range res.Msg.GetBalances() {
		bwe, err := balance.NewWithEntryFromProto(proto)
		if err != nil {
			return nil, err
		}
		bals = append(bals, bwe)
	}
	return bals, nil
}

func (c *client) CreateBalance(ctx context.Context, b *balance.Balance, e *balance.Entry) (*balance.WithEntry, error) {
	req := connect.NewRequest(&balancesv1.CreateRequest{
		TypeId:         b.TypeID,
		CurrencyId:     b.CurrencyID,
		BalancePayload: balance.EncodePayloadToProto(b.Payload),
		EntryYmd:       e.YMD.String(),
		EntryPayload:   balance.EncodePayloadToProto(e.Payload),
	})
	res, err := c.Balances.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	bwe, err := balance.NewWithEntryFromProto(res.Msg.GetBalance())
	if err != nil {
		return nil, err
	}

	return bwe, nil
}
