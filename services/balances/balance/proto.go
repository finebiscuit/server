package balance

import (
	"encoding/base64"

	balancesv1 "github.com/finebiscuit/proto/biscuit/balances/v1"

	"github.com/finebiscuit/server/model/buid"
	"github.com/finebiscuit/server/model/date"
	"github.com/finebiscuit/server/model/payload"
)

func NewWithEntryFromProto(proto *balancesv1.Balance) (*WithEntry, error) {
	id, err := buid.Parse(proto.GetId())
	if err != nil {
		return nil, err
	}
	balancePayload, err := NewPayloadFromProto(proto.GetPayload())
	if err != nil {
		return nil, err
	}
	bwe := &WithEntry{
		Balance: Balance{
			ID:         ID{id},
			TypeID:     proto.GetTypeId(),
			CurrencyID: proto.GetCurrencyId(),
			Payload:    balancePayload,
		},
	}

	if e := proto.GetCurrentEntry(); e != nil {
		ymd, err := date.NewFromString(e.GetYmd())
		if err != nil {
			return nil, err
		}
		entryPayload, err := NewPayloadFromProto(e.GetPayload())
		if err != nil {
			return nil, err
		}
		bwe.Entry = Entry{
			YMD:     ymd,
			Payload: entryPayload,
		}
	}

	return bwe, nil
}

func (b WithEntry) AsProto() *balancesv1.Balance {
	return &balancesv1.Balance{
		Id:         b.ID.String(),
		TypeId:     b.TypeID,
		CurrencyId: b.CurrencyID,
		Payload:    EncodePayloadToProto(b.Payload),
		CurrentEntry: &balancesv1.Entry{
			Ymd:     b.Entry.YMD.String(),
			Payload: EncodePayloadToProto(b.Entry.Payload),
		},
	}
}

func NewPayloadFromProto(p *balancesv1.Payload) (payload.Payload, error) {
	if p == nil {
		return payload.Payload{}, ErrInvalidPayload
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

func EncodePayloadToProto(p payload.Payload) *balancesv1.Payload {
	value := base64.StdEncoding.EncodeToString(p.Blob)
	return &balancesv1.Payload{
		Version:     p.Version,
		Scheme:      uint32(p.Scheme),
		Base64Value: value,
	}
}
