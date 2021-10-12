package custom_err

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"runtime"
	"strings"
)

var (
	ErrDeletingUserNotFound = errors.New("error during deleting user, user was not found")
	ErrDeletingUser         = errors.New("error during deleting user")
	ErrGetUser              = errors.New("error email or password incorrect")
	ErrCreatingUser         = errors.New("error during creating user, email is already exist")
	ErrUnexpected           = errors.New("unexpected")
	ErrPassword             = errors.New("error hashing password")
)

type CustomErrorWithCode struct {
	RestCode int
	GrpcCode codes.Code
	Msg      error
}

func New(originalError error, forUserError error, grpcCode codes.Code) *CustomErrorWithCode {
	return &CustomErrorWithCode{
		Msg:      fmt.Errorf("original error msg: %v, for user error message: %w", originalError, forUserError),
		GrpcCode: grpcCode,
	}
}

func (e *CustomErrorWithCode) Error() string {
	return e.Msg.Error()
}

func (e *CustomErrorWithCode) ErrorForUser() string {
	return UnwrapRecursive(e.Msg).Error()
}

func GetInfo() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	fn := runtime.FuncForPC(pc).Name()
	funcName := fn[strings.LastIndex(fn, ".")+1:]

	return fmt.Sprintf("%v: %%w", funcName)
}

func UnwrapRecursive(err error) error {
	unwrappedError := errors.Unwrap(err)
	if unwrappedError != nil {
		return UnwrapRecursive(unwrappedError)
	}
	return err
}
