package service

//go:generate mockery --name Service

type Service interface {
	User() User
}
