package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type System struct {
	SocksPort   int         `json:"socks_port"`
	LightConfig LightConfig `json:"light_config"`
	Log         LogConfig   `json:"log"`
	UserConfig  string      `json:"user_config"`
}
type LogConfig struct {
	Debug  string `json:"debug"`
	Access string `json:"access"`
}
type LightConfig struct {
	PrivateKey string `json:"private_key"`
	NodeID     string `json:"node_id"`
	Port       int    `json:"port"`
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

func Initialize(path string) (*System, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBytes, &SystemConfig)
	if err != nil {
		return nil, err
	}
	return &SystemConfig, nil
}
