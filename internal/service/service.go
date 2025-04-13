package service

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type IRepository interface {
	Start(ctx context.Context, b *bot.Bot, update *models.Update)
	Support(ctx context.Context, b *bot.Bot, update *models.Update)
	LayoutTARO(ctx context.Context, b *bot.Bot, update *models.Update)
	Numerology(ctx context.Context, b *bot.Bot, update *models.Update)
	Default(ctx context.Context, b *bot.Bot, update *models.Update)
}

type Service struct {
	IRepository IRepository
}

func NewService(IRepository IRepository) *Service {
	return &Service{IRepository: IRepository}
}

func (s *Service) StartService(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.IRepository.Start(ctx, b, update)
}

func (s *Service) SupportService(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.IRepository.Support(ctx, b, update)
}

func (s *Service) LayoutTAROService(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.IRepository.LayoutTARO(ctx, b, update)
}

func (s *Service) NumerologyService(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.IRepository.Numerology(ctx, b, update)
}
func (s *Service) DefaultService(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.IRepository.Default(ctx, b, update)
}
