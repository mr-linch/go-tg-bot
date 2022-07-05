package service

//go:generate mockery --name Service

type Service interface {
	Auth() Auth
}
