package renew

type StatusCode int

const (
	//RUNNING ...
	RUNNING StatusCode = iota
	FETCHING
	NOUPDATEFETCHED
	UPDATEFETCHED
	FAILURE
	RESTARTING
)
