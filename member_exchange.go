package gateway

import (
	"context"
	"strings"

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

	data, err := req.Do(ctx).Bytes()
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
