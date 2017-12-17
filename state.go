package renew

type StatusCode int

const (
	//RUNNING ...
	RUNNING StatusCode = iota
	FETCHING
	UPDATEFETCHED
	FAILURE
)

//State ...
type State struct {
	Description string
	StatusCode  StatusCode
}
