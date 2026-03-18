package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	botToken = "YOUR_BOT_TOKEN_HERE"
	chatID   = 123456789 // Chat ID của Telegram là kiểu số (int64)
)

func main() {
	// 1. Khởi tạo Bot Telegram
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Không thể kết nối Telegram Bot:", err)
	}

	log.Printf("Đã đăng nhập với tài khoản Bot: %s", bot.Self.UserName)

	// 2. Thiết lập Gin
	r := gin.Default()

	r.POST("/notify", func(c *gin.Context) {
		var input struct {
			Message string `json:"message" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tin nhắn không hợp lệ"})
			return
		}

		// 3. Sử dụng thư viện để gửi tin nhắn
		msg := tgbotapi.NewMessage(int64(chatID), input.Message)
		msg.ParseMode = "Markdown" // Hỗ trợ định dạng đậm/nghiêng

		_, err = bot.Send(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gửi Telegram thất bại"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Thành công!"})
	})

	r.Run(":8080")
}