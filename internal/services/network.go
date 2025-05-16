package service

import "resource-graber/internal/domains/models"

func (s *Service) ProcessNetwork() {
	for {
		_, err := s.GetNetwork()
		if err != nil {
			s.logger.Error("error getting network", "error", err)
			continue
		}
	}
}

func (s *Service) GetNetwork() (models.Network, error) {
	return models.Network{}, nil
}
