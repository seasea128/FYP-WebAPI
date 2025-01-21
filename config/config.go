package config

import (
	"sync/atomic"

	"github.com/spf13/viper"
)

func loadDefault() {
	viper.SetDefault("DBType", "sqlite")
	viper.SetDefault("ConnectionString", "test.db")
	viper.SetDefault("HTTPPort", "8080")
	viper.SetDefault("MQTTPort", "1883")
}

type Configuration struct {
	DBType           DBType `json:"DBType"`
	ConnectionString string `json:"ConnectionString"`
	HTTPPort         int    `json:"HTTPPort"`
	MQTTPort         int    `json:"MQTTPort"`
}

type DBType string

const (
	SQLite   DBType = "sqlite"
	Postgres DBType = "postgres"
)

func LoadConfig() (*atomic.Pointer[Configuration], error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	loadDefault()

	config := &Configuration{}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := viper.SafeWriteConfig()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	var atomicCfg atomic.Pointer[Configuration]
	atomicCfg.Store(config)

	return &atomicCfg, nil
}
