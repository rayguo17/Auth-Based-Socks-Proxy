package setup

import (
	"encoding/pem"
	"errors"
	"github.com/rayguo17/go-socks/Backdoor"
	"github.com/rayguo17/go-socks/api"
	"github.com/rayguo17/go-socks/cmd/config"
	"github.com/rayguo17/go-socks/manager"
	"github.com/rayguo17/go-socks/util/logger"
	"gitlab.com/yawning/obfs4.git/common/ntor"
	"log"
	"os"
	"time"
)

func Server(system *config.System) error {
	umStartChan := make(chan bool)
	if system.Interface == "graphic" {
		config := logger.Config{
			AccessPath: system.GetAccessPath(),
			DebugPath:  system.GetDebugPath(),
			IsMulti:    true,
			LogWriter:  system.GetLogWriter(),
		}
		err := logger.InitializeLogger(config)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		config := logger.Config{
			AccessPath: system.GetAccessPath(),
			DebugPath:  system.GetDebugPath(),
			IsMulti:    false,
		}
		err := logger.InitializeLogger(config)
		if err != nil {
			log.Fatal(err)
		}
	}

	privPemByte, err := os.ReadFile(system.LightConfig.PrivateKeyFile)
	if err != nil {
		return err
	}
	prBlock, _ := pem.Decode(privPemByte)
	key := prBlock.Bytes[len(prBlock.Bytes)-32:]
	privKey, err := ntor.NewPublicKey(key)
	//logger.Access.Printf("Server private key: %v\n", privKey.Hex())
	if err != nil {
		return err
	}
	system.LightConfig.PrivateKey = privKey.Hex()
	//pp.Println(system)
	pubPemByte, err := os.ReadFile(system.LightConfig.PublicKeyFile)
	if err != nil {
		return err
	}
	pbBlock, _ := pem.Decode(pubPemByte)
	pubKeyByte := pbBlock.Bytes[len(pbBlock.Bytes)-32:]
	pubKey, err := ntor.NewPublicKey(pubKeyByte)
	system.LightConfig.PublicKey = pubKey.Hex()
	if err != nil {
		return err
	}
	logger.Access.Printf("Server public key: %v\n", pubKey.Hex())
	if err != nil {
		return err
	}
	manager.UM.Initialize(system.GetConfigPath())
	if system.Interface == "graphic" {
		manager.UM.SetSubscribe(system.GetCtx())
	}
	go manager.UM.MainRoutine(umStartChan)
	select {
	case <-umStartChan:
		logger.Access.Println("UM initialize success")
	case <-time.After(time.Second * 5):
		return errors.New("start manager manager timeout")
	}
	if system.IsApiActive() {
		go api.MainRoutine(system.GetApiPort())
	}
	if system.IsBackDoorActive() {
		go Backdoor.BackDoorRoutine()
	}
	//setup logger...
	return nil
}
