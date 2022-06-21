package balances

import (
	"context"
	"errors"
	"net/http"

	"github.com/bufbuild/connect-go"
	balancesv1 "github.com/finebiscuit/proto/biscuit/balances/v1"
	"github.com/finebiscuit/proto/biscuit/balances/v1/balancesv1connect"

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
		Balance: b.AsProto(),
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
		res.Msg.Balances = append(res.Msg.Balances, b.AsProto())
	}
	return res, nil
}

func (h *handler) Create(
	ctx context.Context,
	req *connect.Request[balancesv1.CreateRequest],
) (*connect.Response[balancesv1.CreateResponse], error) {
	bp, err := balance.NewPayloadFromProto(req.Msg.GetBalancePayload())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	b, err := balance.New(req.Msg.GetTypeId(), req.Msg.GetCurrencyId(), bp)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	ep, err := balance.NewPayloadFromProto(req.Msg.GetEntryPayload())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	e, err := balance.NewEntryWithString(req.Msg.GetEntryYmd(), ep)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	bwe, err := h.Balances.CreateBalance(ctx, b, e)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&balancesv1.CreateResponse{
		Balance: bwe.AsProto(),
	})
	return res, nil
}

func (h *handler) UpsertEntry(
	ctx context.Context,
	req *connect.Request[balancesv1.UpsertEntryRequest],
) (*connect.Response[balancesv1.UpsertEntryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
