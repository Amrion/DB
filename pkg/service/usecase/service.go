package usecase

import (
	"forum/pkg/models"
	"forum/pkg/service/repository"
)

type ServiceUsecaseInterface interface {
	CleanDb() bool
	GetStatus() models.Status
}

func (s ServiceUsecase) GetStatus() models.Status {
	return s.DB.GetStatus()
}

func (s ServiceUsecase) CleanDb() bool {
	return s.DB.CleanDb()
}

type ServiceUsecase struct {
	DB repository.ServiceRepository
}
