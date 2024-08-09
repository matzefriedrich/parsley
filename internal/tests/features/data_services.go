package features

type dataService interface {
	FetchData() string
}

type remoteDataService struct {
}

func newRemoteDataService() dataService {
	return &remoteDataService{}
}

func (r *remoteDataService) FetchData() string {
	return "data from remote service"
}

var _ dataService = &remoteDataService{}

type localDataService struct{}

func newLocalDataService() dataService {
	return &localDataService{}
}

func (l *localDataService) FetchData() string {
	return "data from local service"
}

var _ dataService = &localDataService{}
