package main

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var step int
var weight float64
var dosePerKg float64
var concentration float64

func main() {
    bot, err := tgbotapi.NewBotAPI("7535167012:AAHKNv3Z0aD0euf-l5F7E6W_GqD7vIz_yok")
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message != nil {
            // Если сообщение содержит команду /start
            if update.Message.Text == "/start" {
                // Создаём Reply-клавиатуру с кнопками "Рассчитать дозу", "Перезапустить расчёт", "Помощь"
                keyboard := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("Рассчитать дозу"),
                        tgbotapi.NewKeyboardButton("Перезапустить расчёт"),
                    ),
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("Помощь"),
                    ),
                )

                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
                msg.ReplyMarkup = keyboard
                bot.Send(msg)
            }

            // Обработка нажатий на кнопки
            switch update.Message.Text {
            case "Рассчитать дозу":
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите вес животного в кг:")
                bot.Send(msg)
                step = 1
            case "Перезапустить расчёт":
                // Сбрасываем шаги и перезапускаем расчёт
                step = 0
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Расчёт перезапущен. Нажмите 'Рассчитать дозу' для нового расчёта.")
                bot.Send(msg)
            case "Помощь":
                // Отправляем сообщение с инструкциями
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Этот бот поможет вам рассчитать дозу анестезии для животных. Следуйте инструкциям:\n\n1. Нажмите 'Рассчитать дозу'.\n2. Введите вес животного, дозу на кг и концентрацию препарата.\n3. Бот вычислит нужное количество препарата и выведет результат.\n\nВы также можете использовать кнопку 'Перезапустить расчёт' для начала нового расчёта.")
                bot.Send(msg)
                // Сбрасываем шаги
                step = 0
                msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите 'Рассчитать дозу' для начала нового расчёта.")
                bot.Send(msg)
            }

            // Пошаговый процесс расчёта
            if step > 0 {
                switch step {
                case 1:
                    // Получаем вес животного
                    var err error
                    weight, err = strconv.ParseFloat(update.Message.Text, 64)
                    if err != nil {
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите корректное число.")
                        bot.Send(msg)
                    } else {
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите дозу на кг в мг:")
                        bot.Send(msg)
                        step = 2
                    }
                case 2:
                    // Получаем дозу на кг
                    var err error
                    dosePerKg, err = strconv.ParseFloat(update.Message.Text, 64)
                    if err != nil {
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите корректное число.")
                        bot.Send(msg)
                    } else {
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите концентрацию препарата в мг/мл:")
                        bot.Send(msg)
                        step = 3
                    }
                case 3:
                    // Получаем концентрацию и делаем расчет
                    var err error
                    concentration, err = strconv.ParseFloat(update.Message.Text, 64)
                    if err != nil {
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите корректное число.")
                        bot.Send(msg)
                    } else {
                        // Рассчёт дозы
                        dose := weight * dosePerKg
                        volume := dose / concentration

                        // Отправляем результат
                        response := fmt.Sprintf("Для животного весом %.2f кг с дозой %.2f мг/кг и концентрацией %.2f мг/мл нужно %.2f мл препарата.", weight, dosePerKg, concentration, volume)
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
                        bot.Send(msg)

                        // Сбрасываем шаги для нового расчёта
                        step = 0
                    }
                }
            }
        }
    }
}
