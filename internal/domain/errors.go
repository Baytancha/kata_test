package domain

import "errors"

var (
	ErrRecordNotFound = errors.New("record for host not found")
	ErrScan           = errors.New("scan error")
	ErrTx             = errors.New("transaction error")
)
