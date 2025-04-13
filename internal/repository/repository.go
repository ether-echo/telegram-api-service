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

func (r *Repository) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		log.Info("Message nil")
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

func (r *Repository) Support(ctx context.Context, b *bot.Bot, update *models.Update) {

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
		log.Info("Ошибка при отправке сообщения с кнопкой:", err)
	}
}

var check = false

func (r *Repository) LayoutTARO(ctx context.Context, b *bot.Bot, update *models.Update) {
	if check == false {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "Откройте дверь к своему будущему с помощью карт Таро! \n" +
				"Каждая карта Таро - это ключ к вашей судьбе, доверьтесь им и узнайте, как изменить свою жизнь к лучшему. \n\n" +
				"*Задайте вопрос, который вас тревожит и карты подскажут, какие шаги предпринять*.",
			ParseMode: models.ParseModeMarkdownV1,
		})
		if err != nil {
			log.Fatal(err)
		}
		check = true
		return
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Судьба не любит слишком упорных. Я умолкаю до завтрашнего дня.",
		ParseMode: models.ParseModeMarkdownV1,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func (r *Repository) Numerology(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Для лиц, рожденных под числом ??? сегодня ...",
	})
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
