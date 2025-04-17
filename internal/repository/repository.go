package repository

import (
	"context"
	"github.com/ether-echo/telegram-api-service/adapter/producer"
	"github.com/ether-echo/telegram-api-service/internal/model"
	"github.com/ether-echo/telegram-api-service/pkg/config"
	"github.com/ether-echo/telegram-api-service/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	log = logger.Logger().Named("repository").Sugar()
)

type Repository struct {
	BotTG *bot.Bot
	prod  *producer.KafkaProducer
}

func NewRepository(conf *config.Config, prod *producer.KafkaProducer) *Repository {
	b, err := bot.New(conf.Token)
	if err != nil {
		log.Fatal(err)
	}

	commands := []models.BotCommand{
		{Command: "/start", Description: "Перезапустить бота"},
		{Command: "/support", Description: "Поддержка"},
	}

	_, err = b.SetMyCommands(context.Background(), &bot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		log.Fatal("Ошибка при установке команд:", err)
	}

	return &Repository{
		BotTG: b,
		prod:  prod,
	}
}

func (r *Repository) Default(ctx context.Context, b *bot.Bot, update *models.Update) {

	message := model.MessageRequest{
		ChatId:    update.Message.Chat.ID,
		LastName:  update.Message.Chat.LastName,
		FirstName: update.Message.Chat.FirstName,
		Username:  update.Message.Chat.Username,
		Message:   update.Message.Text,
	}

	err := r.prod.SendMessageToKafka(message)
	if err != nil {
		log.Error(err)
	}

}

func (r *Repository) SendMessage(chatId int64, message, url string) {

	var keyboard models.ReplyMarkup

	keyboard = &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "🔮 Расклад ТАРО"},
				{Text: "💸 Нумерология"},
			},
			{
				{Text: "💺 Поддержка"},
			},
		},
		ResizeKeyboard:  true, // Уменьшает клавиатуру
		OneTimeKeyboard: false,
	}

	if url != "" {
		keyboard = &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "Написать в чат",
						URL:  url,
					},
				},
			},
		}
	}

	ctx := context.Background()

	r.BotTG.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        message + "\n" + url,
		ParseMode:   models.ParseModeMarkdownV1,
		ReplyMarkup: keyboard,
	})
}
