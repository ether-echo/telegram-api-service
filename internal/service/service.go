package service

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type IRepository interface {
	Default(ctx context.Context, b *bot.Bot, update *models.Update)
}

type Service struct {
	IRepository IRepository
}

func NewService(IRepository IRepository) *Service {
	return &Service{IRepository: IRepository}
}

func (s *Service) DefaultService(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.IRepository.Default(ctx, b, update)
}
