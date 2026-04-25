package helper

import (
	"errors"

	"github.com/lib/pq"
)

func IsDuplicateError(err error) bool {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}

	return false
}
