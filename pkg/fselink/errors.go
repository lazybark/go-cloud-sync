package fselink

type ErrorCode int

const (
	err_codes_start ErrorCode = iota

	ErrCodeWrongCreds
	ErrForbidden

	err_codes_end
)

func (ec ErrorCode) String() string {
	codes := [...]string{"wrong login or password", "forbidden"}
	if ec <= err_codes_start || ec >= err_codes_end {
		return "unknown error"
	}
	return codes[ec-1]
}

func (ec ErrorCode) Int() int {
	return int(ec)
}
