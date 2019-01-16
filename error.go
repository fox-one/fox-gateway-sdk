package gateway

import (
	"fmt"
)

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Hint string `json:"hint"`
}

func (e Err) Error() string {
	return e.Msg
}

func (e Err) String() string {
	err := fmt.Sprintf("%d : %s", e.Code, e.Msg)
	if len(e.Hint) > 0 {
		err = err + " (" + e.Hint + ")"
	}

	return err
}
