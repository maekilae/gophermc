package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"codeberg.org/makila/minecraftgo/config"
	"codeberg.org/makila/minecraftgo/internal/db"
	"codeberg.org/makila/minecraftgo/internal/logger"
	"codeberg.org/makila/minecraftgo/internal/network"
)

func main() {
	logger.Init("log")
	db, err := db.OpenDB("world", "player")
	if err != nil {
		slog.Error("Could not open db", "Error", err)
		os.Exit(1)
	}
	go commands()
	config.LoadFromPath(context.Background(), "properties.pkl")
	s := network.NewServer("MinecraftServer", ":25565", db)
	s.RunServer()

}

func commands() {
	var cmd string
	for {
		fmt.Scan(&cmd)
	}

}
