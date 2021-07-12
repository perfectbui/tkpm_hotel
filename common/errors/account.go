package errors

import (
	"fmt"
)

var (
	ErrEmailExisted         = fmt.Errorf("EMAIL_EXISTED")
	ErrEmailNotFound        = fmt.Errorf("EMAIL_NOT_FOUND")
	ErrPasswordIsNotCorrect = fmt.Errorf("PASSWORD_INCORRECT")
	ErrInvalidToken         = fmt.Errorf("TOKEN_INVALID")
	ErrNoPermission         = fmt.Errorf("NO_PERMISSION")
)
