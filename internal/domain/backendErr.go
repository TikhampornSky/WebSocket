package domain

import "fmt"

type ErrKind int

const (
	_ ErrKind = iota
	UserEmailNotFound
	UserIDNotFound
	DuplicateEmail
	DuplicateUsername

	DuplicateChatroom
	ChatroomIDNotFound
	ChatroomPrivate
	ChatroomFull
	
	Internal
)

var (
	ErrUserEmailNotFound = BackEndError{Kind: UserEmailNotFound}
	ErrUserIDNotFound    = BackEndError{Kind: UserIDNotFound}
	ErrDuplicateEmail    = BackEndError{Kind: DuplicateEmail}
	ErrDuplicateUsername = BackEndError{Kind: DuplicateUsername}

	ErrDuplicateChatroom = BackEndError{Kind: DuplicateChatroom}
	ErrChatroomIDNotFound  = BackEndError{Kind: ChatroomIDNotFound}
	ErrChatroomPrivate = BackEndError{Kind: ChatroomPrivate}
	ErrChatroomFull = BackEndError{Kind: ChatroomFull}

	ErrInternal = BackEndError{Kind: Internal}
)

type BackEndError struct {
	Kind    ErrKind
	Message string
	Detail  map[string]string
	Err     error
}

func (e BackEndError) Error() string {
	return e.Message
}

func (e BackEndError) Is(err error) bool {
	switch errs := err.(type) {
	case BackEndError:
		return e.Kind == errs.Kind
	default:
		return false
	}
}

func (e BackEndError) With(message string, a ...any) *BackEndError {
	ne := e
	ne.Message = fmt.Sprintf(message, a...)
	return &ne
}

func (e BackEndError) WithDetail(message string, detail map[string]string) *BackEndError {
	ne := e
	ne.Message = message
	ne.Detail = detail
	return &ne
}

func (e BackEndError) From(message string, err error) *BackEndError {
	ne := e
	ne.Message = message
	ne.Err = err
	return &ne
}

func (e BackEndError) FromDetail(message string, detail map[string]string, err error) *BackEndError {
	ne := e
	ne.Message = message
	ne.Detail = detail
	ne.Err = err
	return &ne
}

func (e *BackEndError) Unwrap() error {
	return e.Err
}
