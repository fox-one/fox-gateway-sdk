package gateway

import (
	"context"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func (m *MemberService) ExchangeOrder(ctx context.Context, amount, price, side, symbol, trace, orderType string) (*WalletSnapshotView, error) {
	req := m.POST("/order").Body(map[string]interface{}{
		"amount":   amount,
		"price":    price,
		"side":     side,
		"symbol":   strings.ToUpper(symbol),
		"trace_id": trace,
		"type":     orderType,
	})

	data, err := req.Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
	if err != nil {
		return nil, err
	}

	var resp struct {
		*WalletSnapshotView
		Err
	}

	jsoniter.Unmarshal(data, &resp)
	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.WalletSnapshotView, nil
}

func (m *MemberService) CancelExchangeOrder(ctx context.Context, orderID string) error {
	req := m.DELETE("/order/" + orderID)

	data, err := req.Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
	if err != nil {
		return err
	}

	var resp struct {
		Err
	}

	jsoniter.Unmarshal(data, &resp)
	if resp.Code > 0 {
		return resp.Err
	}

	return nil
}
