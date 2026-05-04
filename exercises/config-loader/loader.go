package configloader

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppName   string `json:"app_name"`
	Version   string `json:"version"`
	DebugMode bool   `json:"debug_mode"`
}

func LoadConfig() (cfg *Config, err error) {
	data, err := os.ReadFile("config-loader/data.json")

	if err != nil {
		fmt.Printf("Error %s", err)
		return nil, err
	}

	var config Config

	erro := json.Unmarshal(data, &config)

	if erro != nil {
		fmt.Printf("Error %s", erro)
		return nil, erro
	}
	return &config, nil

}
