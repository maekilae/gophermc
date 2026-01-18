package main

import (
	"codeberg.org/makila/minecraftgo/internal/logger"
	"codeberg.org/makila/minecraftgo/internal/network"
)

func main() {
	logger.Init("log.json")
	s := network.NewServer("MinecraftServer", "tcp", ":25565")
	s.RunServer()

}
