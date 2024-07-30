package main

import (
	"github.com/joho/godotenv"
	"github.com/vitezprchal/waitlistgo/internal/server"
)

func main() {
	godotenv.Load(".env.example")
	server.InitServer(":8080")
}
