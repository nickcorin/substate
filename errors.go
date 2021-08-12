package substate

import "errors"

var (
	// ErrFunctionReturnArgument is returned when the Substate interface
	// contains a method which returns a function.
	ErrFunctionReturnArgument = errors.New("function return args are not supported")

	// ErrMultipleReturnArguments is returned when the Substate interface
	// contains a method with multiple return arguments.
	ErrMultipleReturnArguments = errors.New("multiple return args are not supported")

	// ErrSubstateNotFound is returned when gensubstate is run in a file which
	// doesn't contain a Substate interface.
	ErrSubstateNotFound = errors.New("no substate interface found")
)
