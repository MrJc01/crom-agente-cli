package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// CLIConfig armazena a configuração da interface do terminal
type CLIConfig struct {
	Theme       string `mapstructure:"theme"`
	LogLevel    string `mapstructure:"log_level"`
	DaemonPort  int    `mapstructure:"daemon_port"`
	AutoConnect bool   `mapstructure:"auto_connect"`
}

var defaultConfig = CLIConfig{
	Theme:       "auto",
	LogLevel:    "info",
	DaemonPort:  9090,
	AutoConnect: true,
}

// Load lê a configuração dos flags, variáveis de ambiente e arquivo yaml (Item 36)
func Load() (*CLIConfig, error) {
	v := viper.New()

	// 1. Configura defaults
	v.SetDefault("theme", defaultConfig.Theme)
	v.SetDefault("log_level", defaultConfig.LogLevel)
	v.SetDefault("daemon_port", defaultConfig.DaemonPort)
	v.SetDefault("auto_connect", defaultConfig.AutoConnect)

	// 2. Configura Environment Variables
	v.SetEnvPrefix("CROM_CLI")
	v.AutomaticEnv() // CROM_CLI_THEME, CROM_CLI_DAEMON_PORT...

	// 3. Configura Arquivo
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(homeDir, ".crom", "cli")
		_ = os.MkdirAll(configDir, 0755)
		v.AddConfigPath(configDir)
		v.SetConfigName("config")
		v.SetConfigType("yaml")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("erro ao ler arquivo de configuracao: %w", err)
			}
			// Se o arquivo nao existe, ignora
		}
	}

	var cfg CLIConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("erro ao converter configuracao: %w", err)
	}

	return &cfg, nil
}
