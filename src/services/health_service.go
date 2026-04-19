package services

type HealthService struct {
	IsReady bool
}

func NewHealthService() *HealthService {
	return &HealthService{IsReady: false}
}

func (s *HealthService) Ready() (bool, error) {
	return s.IsReady, nil
}
