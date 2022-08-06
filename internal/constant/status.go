package constant

type Status int

const (
	StatusNotStarted Status = iota
	StatusInProgress
	StatusOk
	StatusError
)
