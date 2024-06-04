package services

import "github.com/adi-kmt/bitespeed_backend/pkg/repositories"

type Service struct {
	repo *repositories.Repository
}

func NewService(repository *repositories.Repository) *Service {
	return &Service{
		repo: repository,
	}
}
