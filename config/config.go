package config

import (
	"encoding/json"
	"os"
)

type ThumbnailCfg struct {
	StorePath          string `json:"store_path"`
	JpegQuality        int    `json:"jpeg_quality"`
	MaxDimension       int    `json:"max_dimension_px"`
	LifetimeUnusedDays int    `json:"lifetime_unused_days"`
}

type Config struct {
	ListenPort   int          `json:"listen_port"`
	GlobalDBPath string       `json:"global_db_path"`
	Thumbnail    ThumbnailCfg `json:"thumbnail"`
}

var gConfig Config

func Initialize(confFilePath string) error {
	confFile, err := os.Open(confFilePath)
	if err != nil {
		return err
	}
	defer confFile.Close()

	var conf Config
	jsonParser := json.NewDecoder(confFile)
	err = jsonParser.Decode(&conf)
	if err != nil {
		return err
	}

	gConfig = conf
	return nil
}

func Get() *Config {
	return &gConfig
}
