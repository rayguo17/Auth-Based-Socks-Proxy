package config

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
)

//should be able to choose to use encryption or not...
type System struct {
	Interface   string      `json:"interface"`
	SocksPort   int         `json:"socks_port"`
	LightConfig LightConfig `json:"light_config"`
	Log         LogConfig   `json:"log"`
	UserConfig  string      `json:"user_config"`
	BackDoor    BackDoor    `json:"back_door"`
	ApiServer   ApiServer   `json:"api_server"`
	Mode        string      `json:"mode"`
	ctx         context.Context
}
type BackDoor struct {
	active bool `json:"active"`
}
type ApiServer struct {
	active bool `json:"active"`
	port   int  `json:"port"`
}
type LogConfig struct {
	Debug     string `json:"debug"`
	Access    string `json:"access"`
	logWriter io.Writer
}
type LightConfig struct {
	PrivateKeyFile string `json:"private_key_file"`
	PublicKeyFile  string `json:"public_key_file"`
	NodeID         string `json:"node_id"`
	Port           int    `json:"port"`
	PrivateKey     string
	PublicKey      string
}

func (s *System) SetLogWriter(lw io.Writer) {
	s.Log.logWriter = lw
}
func (s *System) GetLogWriter() io.Writer {
	return s.Log.logWriter
}
func (s *System) SetCtx(ctx context.Context) {
	s.ctx = ctx
}
func (s *System) GetCtx() context.Context {
	return s.ctx
}
func (s *System) GetLightPort() string {

	return strconv.Itoa(s.LightConfig.Port)
}

var SystemConfig System

func (s *System) GetSocksPort() string {
	return strconv.Itoa(s.SocksPort)
}
func (s *System) GetConfigPath() string {
	return s.UserConfig
}
func (s *System) GetDebugPath() string {
	return s.Log.Debug
}
func (s *System) GetAccessPath() string {
	return s.Log.Access
}
func (s *System) IsApiActive() bool {
	return s.ApiServer.active
}
func (s *System) GetApiPort() int {
	return s.ApiServer.port
}
func (s *System) IsBackDoorActive() bool {
	return s.BackDoor.active
}

func Initialize(path string) (*System, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBytes, &SystemConfig)
	if err != nil {
		return nil, err
	}
	if SystemConfig.LightConfig.PrivateKeyFile == "" || SystemConfig.LightConfig.PublicKeyFile == "" {
		return nil, errors.New("private key file path should not be null")
	}

	return &SystemConfig, nil
}
