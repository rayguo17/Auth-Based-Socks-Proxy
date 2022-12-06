package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type System struct {
	Port       int       `json:"port"`
	Log        LogConfig `json:"log"`
	UserConfig string    `json:"user_config"`
}
type LogConfig struct {
	Debug  string `json:"debug"`
	Access string `json:"access"`
}

func (s *System) GetPort() string {
	return strconv.Itoa(s.Port)
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
	var system System
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBytes, &system)
	if err != nil {
		return nil, err
	}
	return &system, nil
}
