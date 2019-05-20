package gateway

import (
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Hint string `json:"hint"`
}

func (e Err) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s (%s)", e.Msg, e.Hint)
	}

	return e.Msg
}

func (e Err) String() string {
	err := fmt.Sprintf("%d : %s", e.Code, e.Msg)
	if len(e.Hint) > 0 {
		err = err + " (" + e.Hint + ")"
	}

	return err
}

func decodeErr(data []byte) (e Err) {
	jsoniter.Unmarshal(data, &e)
	return
}

var (
	ErrPinInvalid          = errors.New("invalid pin")          // 1104
	ErrInsufficientBalance = errors.New("insufficient balance") // 1161
)

func gatewayErr(e Err) error {
	switch e.Code {
	case 0:
		return nil
	case 1104:
		return ErrPinInvalid
	case 1161:
		return ErrInsufficientBalance
	default:
		return e
	}
}
