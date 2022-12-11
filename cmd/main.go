package main

import (
	"context"
	"fmt"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/cmd/protocol/light"
	"github.com/rayguo17/go-socks/cmd/protocol/socks"
	"github.com/rayguo17/go-socks/cmd/setup"
	"github.com/rayguo17/go-socks/util/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var config_path = "./config.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error not enough argument")
		return
	}
	conf := os.Args[1]
	fmt.Println("Hello world")
	system, err := config.Initialize(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	//pp.Println(system)
	err = setup.Server(system)
	if err != nil {
		log.Fatal(err.Error())
	}
	logger.Access.Printf("Socks server listening at port %v\n", system.GetSocksPort())
	logger.Access.Printf("light server listening at port %v\n", system.GetLightPort())
	socksCtx, cancelFunc := context.WithCancel(context.Background())
	err = socks.ListenStart(system, cancelFunc)
	if err != nil {
		logger.Debug.Fatal(err)
	}
	lightCtx, cancelFunc := context.WithCancel(context.Background())
	err = light.ListenStart(system, cancelFunc)
	if err != nil {
		logger.Debug.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-lightCtx.Done():
		logger.Debug.Println("light terminated")
	case <-socksCtx.Done():
		logger.Debug.Println("socks light terminated")
	case <-sigs:
		logger.Debug.Println("SIGINT received ending")
	}

}
