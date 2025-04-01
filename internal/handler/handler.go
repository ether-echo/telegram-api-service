package handler

import (
	"context"
	"github.com/ether-echo/telegram-api-service/pkg/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

type Handler struct {
	botTG *bot.Bot
}

func NewHandler(ctx context.Context, conf *config.Config) *Handler {

	b, err := bot.New(conf.Token)
	if err != nil {
		panic(err)
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

	return &Handler{
		botTG: b,
	}
}

func (h *Handler) RegisterHandler(ctx context.Context) {

	h.botTG.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, h.startHandler)

	h.botTG.RegisterHandler(bot.HandlerTypeMessageText, "/support", bot.MatchTypePrefix, h.supportHandler)
	h.botTG.RegisterHandler(bot.HandlerTypeMessageText, "💺 Поддержка", bot.MatchTypePrefix, h.supportHandler)

	h.botTG.Start(ctx)
}

func (h *Handler) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		log.Println("Message nil")
		return
	}

	keyboard := &models.ReplyKeyboardMarkup{
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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "🔮Приветствую тебя в мире ТАРО.\n\n *Хочешь найти ответы на интересующие тебя вопросы?*" +
			"\n\nМои карты помогут найти ответы по направлениям:\n" +
			"*- конкретный вопрос*\n" +
			"*- день*\n" +
			"*- любовь*\n" +
			"*- успех*\n" +
			"*- будущий год*\n" +
			"*Поддержка* /support\n" +
			"Нужна помощь? Напишите нам в чат\n\n" +
			"Я стану твоим проводником по лабиринтам судьбы, *помогу разгадать тайные значения карт* и поделюсь мудростью," +
			" которую таит в себе будущее.\n" +
			"💫 Выбери в меню на что будем делать расклад и следуй подсказкам👇",
		ParseMode:   models.ParseModeMarkdownV1,
		ReplyMarkup: keyboard,
	})
}

func (h *Handler) supportHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "Написать в чат",
					URL:  "https://t.me/Degprt",
				},
			},
		},
	}

	messageText := "Нужна помощь? Напишите мне в чат!"

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        messageText,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		log.Println("Ошибка при отправке сообщения с кнопкой:", err)
	}

	//linkOptions := &models.LinkPreviewOptions{
	//	URL:           &url,
	//	ShowAboveText: &showText,
	//}

	//b.SendMessage(ctx, &bot.SendMessageParams{
	//	ChatID: update.Message.Chat.ID,
	//	Text: "Нужна помощь? Напишите мне в чат!\n" +
	//		"Кнопка 1 - http://example1.com\n" +
	//		"Кнопка 2 - http://example2.com",
	//	ReplyMarkup: keyboard,
	//	ParseMode:          models.ParseModeMarkdownV1,
	//})
}
