package fetcher

//IFetcher interface
type IFetcher interface {
	Perform()
	ShouldRun() bool
	Init()
}

//Perform update check
func Perform(f IFetcher) {
	f.Perform()
}

//ShouldRun ...
func ShouldRun(f IFetcher) bool {
	return f.ShouldRun()
}

//Init ...
func Init(f IFetcher) {
	f.Init()
}
