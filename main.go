package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/rayguo17/go-socks/app"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/cmd/protocol/light"
	"github.com/rayguo17/go-socks/cmd/protocol/socks"
	"github.com/rayguo17/go-socks/cmd/setup"
	"github.com/rayguo17/go-socks/util/logger"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//go:embed all:frontend/dist
var assets embed.FS
var config_path = "./config.json"

func main() {
	system, err := config.Initialize(config_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	//setup frontend
	//pp.Println(system)
	err = setup.Server(system)
	if err != nil {
		log.Fatal(err.Error())
	}
	logger.Access.Printf("Socks server listening at port %v\n", system.GetSocksPort())
	logger.Access.Printf("light server listening at port %v\n", system.GetLightPort())
	socksCtx, socksCancelFunc := context.WithCancel(context.Background())
	err = socks.ListenStart(system, socksCancelFunc)
	if err != nil {
		logger.Debug.Fatal(err)
	}
	lightCtx, lightCancelFunc := context.WithCancel(context.Background())
	err = light.ListenStart(system, lightCancelFunc)
	if err != nil {
		logger.Debug.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	renderCtx, renderCancelFunc := context.WithCancel(context.Background())
	go RenderMainRoutine(renderCancelFunc)
	select {
	case <-lightCtx.Done():
		logger.Debug.Println("light terminated")
	case <-socksCtx.Done():
		logger.Debug.Println("socks light terminated")
	case <-sigs:
		logger.Debug.Println("SIGINT received ending")
	case <-renderCtx.Done():
		logger.Debug.Println("Render terminated")
		//should handle end.
		return
	}
}

func RenderMainRoutine(cancelFunc context.CancelFunc) {
	app := app.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "myproject",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
	cancelFunc()
	return
}
