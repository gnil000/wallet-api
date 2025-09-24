package handlers

import (
	"errors"
	"net/http"
	"strings"
	"wallet-api/src/models"

	"github.com/google/uuid"
)

var (
	ErrPathIsEmpty = errors.New("path is empty")
)

func validateOperationType(op string) bool {
	if op == models.Operation_type_deposit || op == models.Operation_type_withdraw {
		return true
	}
	return false
}

func validateAmount(amount int64) bool {
	return amount > 0
}

func extractIdFromPath(r *http.Request) (uuid.UUID, error) {
	path := r.URL.Path
	partsOfPath := strings.Split(path, "/")
	if len(partsOfPath) == 0 {
		return uuid.Nil, ErrPathIsEmpty
	}
	id, err := uuid.Parse(partsOfPath[len(partsOfPath)-1])
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
