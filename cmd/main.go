package main

import (
	"codeberg.org/makila/minecraftgo/internal/server"
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":25565")
	if err != nil {
		panic(err)
	}
	fmt.Println("Minecraft Server List Ping Server running on :25565")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go server.HandleConnection(conn)
	}
}

