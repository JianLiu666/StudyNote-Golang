package e

type CODE int

const (
	SUCCESS CODE = iota
	INVALID_PARAMS
	ERROR
	ERROR_ADD_DUPLICATED_USER
	ERROR_USER_NOT_FOUND
)

type GENDER = int

const (
	GIRL GENDER = iota
	BOY
)
