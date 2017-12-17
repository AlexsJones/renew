package fetcher

//IFetcher interface
type IFetcher interface {
	Perform() error
	ShouldRun() bool
	Init()
}

//Perform update check
func Perform(f IFetcher) error {
	return f.Perform()
}

//ShouldRun ...
func ShouldRun(f IFetcher) bool {
	return f.ShouldRun()
}

//Init ...
func Init(f IFetcher) {
	f.Init()
}
