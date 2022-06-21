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

func (h *handler) GetBalance(
	ctx context.Context,
	req *connect.Request[balancesv1.GetBalanceRequest],
) (*connect.Response[balancesv1.GetBalanceResponse], error) {
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

	res := connect.NewResponse(&balancesv1.GetBalanceResponse{
		Balance: b.AsProto(),
	})
	return res, nil
}

func (h *handler) ListBalances(
	ctx context.Context,
	req *connect.Request[balancesv1.ListBalancesRequest],
) (*connect.Response[balancesv1.ListBalancesResponse], error) {
	bals, err := h.Balances.ListBalances(ctx, balance.Filter{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&balancesv1.ListBalancesResponse{
		Balances: make([]*balancesv1.Balance, 0, len(bals)),
	})

	for _, b := range bals {
		res.Msg.Balances = append(res.Msg.Balances, b.AsProto())
	}
	return res, nil
}

func (h *handler) CreateBalance(
	ctx context.Context,
	req *connect.Request[balancesv1.CreateBalanceRequest],
) (*connect.Response[balancesv1.CreateBalanceResponse], error) {
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

	res := connect.NewResponse(&balancesv1.CreateBalanceResponse{
		Balance: bwe.AsProto(),
	})
	return res, nil
}

func (h *handler) UpdateBalance(
	ctx context.Context,
	req *connect.Request[balancesv1.UpdateBalanceRequest],
) (*connect.Response[balancesv1.UpdateBalanceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (h *handler) GetEntry(
	ctx context.Context,
	req *connect.Request[balancesv1.GetEntryRequest],
) (*connect.Response[balancesv1.GetEntryResponse], error) {
	// TODO implement GetEntry
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (h *handler) ListEntries(
	ctx context.Context,
	req *connect.Request[balancesv1.ListEntriesRequest],
) (*connect.Response[balancesv1.ListEntriesResponse], error) {
	// TODO implement ListEntries
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (h *handler) CreateEntry(
	ctx context.Context,
	req *connect.Request[balancesv1.CreateEntryRequest],
) (*connect.Response[balancesv1.CreateEntryResponse], error) {
	balanceID, err := balance.ParseID(req.Msg.GetBalanceId())
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

	if err := h.Balances.CreateEntry(ctx, balanceID, e); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&balancesv1.CreateEntryResponse{
		Entry: &balancesv1.Entry{
			Ymd:     e.YMD.String(),
			Payload: balance.EncodePayloadToProto(e.Payload),
		},
	})
	return res, nil
}

func (h *handler) UpdateEntry(
	ctx context.Context,
	req *connect.Request[balancesv1.UpdateEntryRequest],
) (*connect.Response[balancesv1.UpdateEntryResponse], error) {
	balanceID, err := balance.ParseID(req.Msg.GetBalanceId())
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

	if err := h.Balances.UpdateEntry(ctx, balanceID, e, req.Msg.GetVersionMatch()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&balancesv1.UpdateEntryResponse{
		Entry: &balancesv1.Entry{
			Ymd:     e.YMD.String(),
			Payload: balance.EncodePayloadToProto(e.Payload),
		},
	})
	return res, nil
}
