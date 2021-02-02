package main

import (
	"fmt"

	echo2 "github.com/tiger-game/echo/server/echo"

	"github.com/tiger-game/tiger/xserver"
)

func main() {
	s := xserver.NewServer(echo2.NewServer(), xserver.ServerConfig{
		Frame:  0,
		IP:     "127.0.0.1",
		Port:   2233,
		LogDir: "./log",
	})
	if s == nil {
		fmt.Printf("New Server fail.\n")
		return
	}

	if err := s.Init(); err != nil {
		fmt.Printf("Init Server error: %v\n", err)
		return
	}

	if err := s.Run(); err != nil {
		fmt.Printf("Run Server error: %v\n", err)
	}
}
