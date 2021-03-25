package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	appName                    = "turl.to"
	configuration              *Configuration
	configurationFileName      = "config"
	configurationFileExtension = ".yml"
	configurationFileType      = "yaml"

	storeDirectory           = "./store/"
	configurationFileAbsPath = filepath.Join(storeDirectory, configurationFileName)
	storageFileAbsPath       = filepath.Join(storeDirectory, "db.txt")
)

// Configuration is the root configuration
type Configuration struct {
	Server  ServerConfiguration
	Storage StorageConfiguration
}

// ServerConfiguration setup for the server
type ServerConfiguration struct {
	Port    string
	Hash    string
	Timeout int
}

// StorageConfiguration setup for the storage
type StorageConfiguration struct {
	File string
}

// SetupConfigurationDefaults setup for root configuration
func SetupConfigurationDefaults() (*Configuration, error) {
	viper.SetConfigName(configurationFileName)
	viper.SetConfigType(configurationFileType)
	viper.AddConfigPath(storeDirectory)

	bindEnv()
	setDefaults()

	if err := readConfiguration(); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}

// readConfiguration from file
func readConfiguration() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, err := os.Stat(configurationFileAbsPath + configurationFileExtension); os.IsNotExist(err) {
			os.Create(configurationFileAbsPath + configurationFileExtension)
		} else {
			return err
		}

		// let's write defaults
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}

	return nil
}

func bindEnv() {
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.hash", "TT_SERVER_HASH")
	viper.BindEnv("server.timeout", "TT_SERVER_TIMEOUT")

	viper.BindEnv("storage.file", "STORAGE_FILE")
}

func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.hash", "HASH_KEY")
	viper.SetDefault("server.timeout", 24)

	viper.SetDefault("storage.file", storageFileAbsPath)
}
