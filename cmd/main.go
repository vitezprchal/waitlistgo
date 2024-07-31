package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/vitezprchal/waitlistgo/internal/server"
)

func main() {
	godotenv.Load(".env")
	server.InitServer(os.Getenv("PORT_SERVER"))
}
