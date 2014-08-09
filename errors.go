package rest

import (
	"errors"
)

var (
	// ErrInvalidPrefix is returned when attemping to create a New() client with
	// an invalid prefix.
	ErrInvalidPrefix = errors.New(`Expecting a valid canonical URL (http://...).`)

	// ErrCouldNotCreateMultipart is returned when attemping to create a
	// multipart request with a nil body.
	ErrCouldNotCreateMultipart = errors.New(`Couldn't create a multipart request without a body.`)

	// ErrCouldNotConvert is returned when the request response can't be
	// converted to the expected datatype.
	ErrCouldNotConvert = errors.New(`Could not convert response %s to %s.`)

	// ErrDestinationNotAPointer is returned when attemping to provide a
	// destination that is not a pointer.
	ErrDestinationNotAPointer = errors.New(`Destination is not a pointer.`)
)
