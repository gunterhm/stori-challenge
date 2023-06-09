package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
	"strings"
)

// Configurations Application wide configurations
type Configurations struct {
	Smtp          SmptConfigurations          `koanf:"smtp"`
	Database      DatabaseConfigurations      `koanf:"database"`
	TxnProcessing TxnProcessingConfigurations `koanf:"txnprocessing"`
	Schedule      ScheduleConfigurations      `koanf:"schedule"`
}

// SmptConfigurations SMPT configurations
type SmptConfigurations struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	From     string `koanf:"from"`
	SendMail bool   `koanf:"sendmail"`
}

// DatabaseConfigurations Database configurations
type DatabaseConfigurations struct {
	Dsn string `koanf:"dsn"`
}

// TxnProcessingConfigurations Transaction Processing configurations
type TxnProcessingConfigurations struct {
	IncomingDir    string `koanf:"incomingdir"`
	ArchiveDir     string `koanf:"archivedir"`
	FileNameRegExp string `koanf:"filenameregexp"`
}

// ScheduleConfigurations Schedule configurations
type ScheduleConfigurations struct {
	Seconds int `koanf:"seconds"`
}

// LoadConfig Loads configurations depending upon the environment
func LoadConfig(logger *zap.SugaredLogger) *Configurations {
	k := koanf.New(".")
	err := k.Load(file.Provider("resources/config.yml"), yaml.Parser())
	if err != nil {
		logger.Fatalf("Failed to locate configurations. %v", err)
	}

	// Searches for env variables and will transform them into koanf format
	// e.g. SERVER_PORT variable will be server.port: value
	err = k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil)
	if err != nil {
		logger.Fatalf("Failed to replace environment variables. %v", err)
	}

	var configuration Configurations

	err = k.Unmarshal("", &configuration)
	if err != nil {
		logger.Fatalf("Failed to load configurations. %v", err)
	}

	return &configuration
}
