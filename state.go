package renew

type StatusCode int

const (
	//RUNNING ...
	RUNNING StatusCode = iota
	FETCHING
	UPDATEFETCHED
	FAILURE
	RESTARTING
)
