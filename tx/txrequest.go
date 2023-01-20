package txrequest

type RequestStyle int

const (
	RequestNone     = RequestStyle(0)
	RequestReadonly = RequestStyle(1)
	RequestWrite    = RequestStyle(2)
)
