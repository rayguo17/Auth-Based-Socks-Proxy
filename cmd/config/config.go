package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type System struct {
	SocksPort  int       `json:"socks_port"`
	LightPort  int       `json:"light_port"`
	Log        LogConfig `json:"log"`
	UserConfig string    `json:"user_config"`
}
type LogConfig struct {
	Debug  string `json:"debug"`
	Access string `json:"access"`
}

func (s *System) GetLightPort() string {
	return strconv.Itoa(s.LightPort)
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
