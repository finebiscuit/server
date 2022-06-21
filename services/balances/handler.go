package balances

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/bufbuild/connect-go"
	balancesv1 "github.com/finebiscuit/proto/biscuit/balances/v1"
	"github.com/finebiscuit/proto/biscuit/balances/v1/balancesv1connect"

	"github.com/finebiscuit/server/model/payload"
	"github.com/finebiscuit/server/services/balances/balance"
)

// NewHandler builds an HTTP handler for the Balances service. It returns the path on which to mount
// the handler and the handler itself.
func NewHandler(service Service, opts ...connect.HandlerOption) (string, http.Handler) {
	h := &handler{Balances: service}
	return balancesv1connect.NewBalancesHandler(h, opts...)
}

type handler struct {
	Balances Service
}

var _ balancesv1connect.BalancesHandler = &handler{}

func (h *handler) Get(
	ctx context.Context,
	req *connect.Request[balancesv1.GetRequest],
) (*connect.Response[balancesv1.GetResponse], error) {
	id, err := balance.ParseID(req.Msg.GetBalanceId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	b, err := h.Balances.GetBalance(ctx, id)
	if err != nil {
		if errors.Is(err, balance.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&balancesv1.GetResponse{
		Balance: balanceToProto(b),
	})
	return res, nil
}

func (h *handler) List(
	ctx context.Context,
	req *connect.Request[balancesv1.ListRequest],
) (*connect.Response[balancesv1.ListResponse], error) {
	bals, err := h.Balances.ListBalances(ctx, balance.Filter{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&balancesv1.ListResponse{
		Balances: make([]*balancesv1.Balance, 0, len(bals)),
	})

	for _, b := range bals {
		res.Msg.Balances = append(res.Msg.Balances, balanceToProto(b))
	}
	return res, nil
}

func (h *handler) Create(
	ctx context.Context,
	req *connect.Request[balancesv1.CreateRequest],
) (*connect.Response[balancesv1.CreateResponse], error) {
	bp, err := protoToPayload(req.Msg.GetBalancePayload())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	b, err := balance.New(req.Msg.GetTypeId(), req.Msg.GetCurrencyId(), bp)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	ep, err := protoToPayload(req.Msg.GetEntryPayload())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	e, err := balance.NewEntryWithString(req.Msg.GetEntryYmd(), ep)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := h.Balances.CreateBalance(ctx, b, e); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	bwe, err := h.Balances.GetBalance(ctx, b.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&balancesv1.CreateResponse{
		Balance: balanceToProto(bwe),
	})
	return res, nil
}

func (h *handler) UpsertEntry(
	ctx context.Context,
	req *connect.Request[balancesv1.UpsertEntryRequest],
) (*connect.Response[balancesv1.UpsertEntryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func balanceToProto(b *balance.WithEntry) *balancesv1.Balance {
	return &balancesv1.Balance{
		Id:         b.ID.String(),
		TypeId:     b.TypeID,
		CurrencyId: b.CurrencyID,
		Payload:    payloadToProto(b.Payload),
		CurrentEntry: &balancesv1.Entry{
			Ymd:     b.Entry.YMD.String(),
			Payload: payloadToProto(b.Entry.Payload),
		},
	}
}

func payloadToProto(p payload.Payload) *balancesv1.Payload {
	value := base64.StdEncoding.EncodeToString(p.Blob)
	return &balancesv1.Payload{
		Version:     p.Version,
		Scheme:      uint32(p.Scheme),
		Base64Value: value,
	}
}

func protoToPayload(p *balancesv1.Payload) (payload.Payload, error) {
	if p == nil {
		return payload.Payload{}, errors.New("payload is empty")
	}
	s, err := payload.NewScheme(int(p.GetScheme()))
	if err != nil {
		return payload.Payload{}, err
	}
	blob, err := base64.StdEncoding.DecodeString(p.GetBase64Value())
	if err != nil {
		return payload.Payload{}, err
	}
	return payload.New(s, p.Version, blob)
}
