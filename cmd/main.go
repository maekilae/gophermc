package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/maekilae/gophermc/config"
	db "github.com/maekilae/gophermc/internal/database"
	"github.com/maekilae/gophermc/internal/logger"
	"github.com/maekilae/gophermc/internal/server"
)

func main() {
	logger.Init("log")
	db, err := db.OpenDB("world", "player")
	if err != nil {
		slog.Error("Could not open db", "Error", err)
		os.Exit(1)
	}
	go commands()
	c, err := config.LoadFromPath(context.Background(), "properties.pkl")
	if err != nil {
		slog.Error("Could not load config", "Error", err)
		os.Exit(1)
	}
	s := server.NewServer("MinecraftServer", ":25565", db, c.Version, c.Properties)
	s.RunServer()

}

func commands() {
	var cmd string
	for {
		fmt.Scan(&cmd)
	}

}
