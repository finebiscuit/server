package forex

import (
	"context"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	forexv1 "github.com/finebiscuit/proto/biscuit/forex/v1"
	"github.com/finebiscuit/proto/biscuit/forex/v1/forexv1connect"
)

func NewHandler(service Service, opts ...connect.HandlerOption) (string, http.Handler) {
	h := &handler{Forex: service}
	return forexv1connect.NewForexHandler(h, opts...)
}

type handler struct {
	Forex Service
}

func (h handler) GetRate(
	ctx context.Context,
	req *connect.Request[forexv1.GetRateRequest],
) (*connect.Response[forexv1.GetRateResponse], error) {
	date := time.Now()
	if req.Msg.GetHistorical() {
		date = req.Msg.GetTimestamp().AsTime()
	}
	rate, err := h.Forex.GetRate(ctx, req.Msg.GetBase(), req.Msg.GetTarget(), date)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	res := connect.NewResponse(&forexv1.GetRateResponse{
		Value: rate.String(),
	})
	return res, nil
}

func (h handler) ListRates(
	ctx context.Context, req *connect.Request[forexv1.ListRatesRequest],
) (*connect.Response[forexv1.ListRatesResponse], error) {
	date := time.Now()
	if req.Msg.GetHistorical() {
		date = req.Msg.GetTimestamp().AsTime()
	}

	rates, err := h.Forex.ListRates(ctx, req.Msg.GetBase(), req.Msg.GetTargets(), date)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	res := connect.NewResponse(&forexv1.ListRatesResponse{
		Base:   req.Msg.GetBase(),
		Values: make(map[string]string),
	})
	for cur, rate := range rates {
		res.Msg.Values[cur] = rate.String()
	}

	return res, nil
}
