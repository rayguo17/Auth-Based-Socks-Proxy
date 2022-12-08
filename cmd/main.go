package main

import (
	"context"
	"fmt"
	"github.com/rayguo17/go-socks/cmd/config"
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
	fmt.Println("Hello world")
	system, err := config.Initialize(config_path)
	//pp.Println(system)
	err = setup.Server(system)
	if err != nil {
		log.Fatal(err.Error())
	}
	logger.Access.Printf("Socks light listening at port %v\n", system.GetSocksPort())
	socksCtx, cancelFunc := context.WithCancel(context.Background())
	socks.ListenStart(system, cancelFunc)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-socksCtx.Done():
		logger.Debug.Println("socks light terminated")
	case <-sigs:
		logger.Debug.Println("SIGINT received ending")
	}

}
