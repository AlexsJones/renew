package fetcher

//IFetcher interface
type IFetcher interface {
	Perform(applicationBasePath string) error
	ShouldRun() bool
	Init()
}

//Perform update check
func Perform(f IFetcher, applicationBasePath string) error {
	return f.Perform(applicationBasePath)
}

//ShouldRun ...
func ShouldRun(f IFetcher) bool {
	return f.ShouldRun()
}

//Init ...
func Init(f IFetcher) {
	f.Init()
}
