package setup

import (
	"errors"
	"github.com/rayguo17/go-socks/Backdoor"
	"github.com/rayguo17/go-socks/api"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/user"
	"github.com/rayguo17/go-socks/util/logger"
	"time"
)

func Server(system *config.System) error {
	umStartChan := make(chan bool)
	err := logger.InitializeLogger(system.GetAccessPath(), system.GetDebugPath())
	if err != nil {
		return err
	}
	user.UM.Initialize(system.GetConfigPath())

	go user.UM.MainRoutine(umStartChan)
	select {
	case <-umStartChan:
		logger.Access.Println("UM initialize success")
	case <-time.After(time.Second * 5):
		return errors.New("start user manager timeout")
	}
	go api.MainRoutine()
	go Backdoor.BackDoorRoutine()

	//setup logger...
	return nil
}
