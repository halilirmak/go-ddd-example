package errors

type ErrorType struct {
	t string
}

var (
	ErrorUnknown      = ErrorType{"unknown"}
	ErrorInvalidInput = ErrorType{"invalid-input"}
	ErrorNotFound     = ErrorType{"not-found"}
)

type DetailedError interface {
	Error() string
	Context() string
	ErrorType() ErrorType
	Wrap(er error) error
	UnWrap() error
}

type ContextualError struct {
	error        string
	context      string
	errorType    ErrorType
	wrappedError error
}

func (c ContextualError) Error() string {
	return c.error
}

func (c ContextualError) Context() string {
	return c.context
}

func (c ContextualError) ErrorType() ErrorType {
	return c.errorType
}

func (c *ContextualError) Wrap(err error) error {
	c.wrappedError = err
	return c
}

func (c *ContextualError) UnWrap() error {
	return c.wrappedError
}

func NewContextualError(error, context string) DetailedError {
	return &ContextualError{
		error:     error,
		context:   context,
		errorType: ErrorUnknown,
	}
}

func NewInvalidInputError(error, context string) DetailedError {
	return &ContextualError{
		error:     error,
		context:   context,
		errorType: ErrorInvalidInput,
	}
}

func NewNotFoundError(error, context string) DetailedError {
	return &ContextualError{
		error:     error,
		context:   context,
		errorType: ErrorNotFound,
	}
}
