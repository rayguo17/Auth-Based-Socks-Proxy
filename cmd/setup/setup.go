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
	"os"
	"time"
)

func Server(system *config.System) error {
	umStartChan := make(chan bool)
	err := logger.InitializeLogger(system.GetAccessPath(), system.GetDebugPath())
	privPemByte, err := os.ReadFile(system.LightConfig.PrivateKeyFile)
	if err != nil {
		return err
	}
	prBlock, _ := pem.Decode(privPemByte)
	key := prBlock.Bytes[len(prBlock.Bytes)-32:]
	privKey, err := ntor.NewPublicKey(key)
	if err != nil {
		return err
	}
	system.LightConfig.PrivateKey = privKey.Hex()
	pubPemByte, err := os.ReadFile(system.LightConfig.PrivateKeyFile)
	if err != nil {
		return err
	}
	pbBlock, _ := pem.Decode(pubPemByte)
	pubKeyByte := pbBlock.Bytes[len(pbBlock.Bytes)-32:]
	pubKey, err := ntor.NewPublicKey(pubKeyByte)
	if err != nil {
		return err
	}
	logger.Access.Printf("Server public key: %v\n", pubKey.Hex())
	if err != nil {
		return err
	}
	manager.UM.Initialize(system.GetConfigPath())

	go manager.UM.MainRoutine(umStartChan)
	select {
	case <-umStartChan:
		logger.Access.Println("UM initialize success")
	case <-time.After(time.Second * 5):
		return errors.New("start manager manager timeout")
	}
	go api.MainRoutine()
	go Backdoor.BackDoorRoutine()

	//setup logger...
	return nil
}
