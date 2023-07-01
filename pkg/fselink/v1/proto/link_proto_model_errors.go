package proto

import (
	"encoding/json"
	"fmt"
)

type MessageError struct {
	Error     string
	ErrorCode int
}

func (em ExchangeMessage) ReadError() (MessageError, error) {
	var m MessageError
	if em.Type != MessageTypeError {
		return m, fmt.Errorf("[ExchangeMessage][ReadError] unexpected message type '%s'", em.Type)
	}

	err := json.Unmarshal(em.Payload, &m)
	if err != nil {
		return m, fmt.Errorf("[ExchangeMessage][ReadError] %w", err)
	}

	return m, nil
}

type ErrorCode int

var errCodes = [...]string{
	"wrong login or password",
	"forbidden",
	"wrong object type",
	"reading filed",
	"internal server error",
	"message reading failed",
	"object does not exist",
	"unexpected message type recieved",
}

const (
	err_codes_start ErrorCode = iota

	ErrCodeWrongCreds
	ErrForbidden
	ErrWrongObjectType
	ErrFileReadingFailed
	ErrInternalServerError
	ErrMessageReadingFailed
	ErrNotExists
	ErrUnexpectedMessageType

	err_codes_end
)

func (ec ErrorCode) String() string {

	if ec <= err_codes_start || ec >= err_codes_end {
		return "unknown error"
	}
	return errCodes[ec-1]
}

func (ec ErrorCode) Int() int {
	return int(ec)
}
