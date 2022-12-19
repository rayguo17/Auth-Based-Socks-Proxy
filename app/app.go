package app

import (
	"context"
	"fmt"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/cmd/protocol/light"
	"github.com/rayguo17/go-socks/cmd/protocol/socks"
	"github.com/rayguo17/go-socks/cmd/setup"
	"github.com/rayguo17/go-socks/manager"
	"github.com/rayguo17/go-socks/util/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// App struct
type App struct {
	ctx context.Context
	log *LogWriter
}

var config_path = "./desktop_config.json"

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.log = NewLogger(ctx)
	system, err := config.Initialize(config_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	system.SetCtx(ctx)
	system.SetLogWriter(a.log)
	//setup frontend
	//pp.Println(system)
	err = setup.Server(system)
	if err != nil {
		log.Fatal(err.Error())
	}
	logger.Access.Printf("Socks server listening at port %v\n", system.GetSocksPort())
	logger.Access.Printf("light server listening at port %v\n", system.GetLightPort())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	if system.Mode == "debug" {
		logger.Access.Println("Running in debug mode!")

	} else {
		_, socksCancelFunc := context.WithCancel(context.Background())
		err = socks.ListenStart(system, socksCancelFunc)
		if err != nil {
			logger.Debug.Fatal(err)
		}
		_, lightCancelFunc := context.WithCancel(context.Background())
		err = light.ListenStart(system, lightCancelFunc)
		if err != nil {
			logger.Debug.Fatal(err)
		}
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) GetConfig() *Response {
	res := ConfigHandler()
	return res
}
func (a *App) ListUser() []*manager.Display {
	resp := GetAllUser()
	if resp.GetErrCode() != 0 {
		return []*manager.Display{}
	} else {
		return resp.GetData().([]*manager.Display)
	}
}

func (a *App) AddUser(newUser interface{}) *Response {
	//pp.Println(newUser)
	res := AddUserHandler(newUser)
	if res.GetErrCode() != 0 {
		logger.Debug.Println(res.GetErrMsg())
	}
	return NewAppResponse(res.GetErrCode(), res.GetErrMsg(), nil)
	//pp.Println(u)
}

func (a *App) DelUser(name string) *Response {
	res := DelUserHandler(name)
	if res.GetErrCode() != 0 {
		logger.Debug.Println(res.GetErrMsg())
	}
	return NewAppResponse(res.GetErrCode(), res.GetErrMsg(), nil)
}
