package handler

import (
	"context"
	"github.com/ether-echo/telegram-api-service/internal/repository"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type IService interface {
	DefaultService(ctx context.Context, b *bot.Bot, update *models.Update)
}

type Handler struct {
	IService IService
}

func NewHandler(service IService) *Handler {
	return &Handler{
		IService: service,
	}
}

func (h *Handler) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.IService.DefaultService(ctx, b, update)
}

func (h *Handler) StartBot(ctx context.Context, rep *repository.Repository) {

	//rep.BotTG.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, h.StartHandler)
	//rep.BotTG.RegisterHandler(bot.HandlerTypeMessageText, "/support", bot.MatchTypePrefix, h.SupportHandler)
	//rep.BotTG.RegisterHandler(bot.HandlerTypeMessageText, "💺 Поддержка", bot.MatchTypePrefix, h.SupportHandler)
	//rep.BotTG.RegisterHandler(bot.HandlerTypeMessageText, "🔮 Расклад ТАРО", bot.MatchTypePrefix, h.LayoutTAROHandler)
	//rep.BotTG.RegisterHandler(bot.HandlerTypeMessageText, "💸 Нумерология", bot.MatchTypePrefix, h.NumerologyHandler)
	rep.BotTG.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, h.DefaultHandler)

	rep.BotTG.Start(ctx)
}
