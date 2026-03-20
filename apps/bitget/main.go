package main

import (
	"log"
)

func main() {
	// config.Env()
	// // database.Connect()
	// // database.Migrate()
	r := routes.Routes()
	// appPort := config.AppPort()


	port := "8080" // có thể lấy từ config sau
	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed:", err)
	}
}