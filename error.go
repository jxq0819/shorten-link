package main

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (statusError StatusError) Status() int {
	return statusError.Code
}

func (statusError StatusError) Error() string {
	return statusError.Err.Error()
}
