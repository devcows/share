package lib

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

type SettingsShare struct {
	Daemon Daemon
}

type Daemon struct {
	Port             int
	Host             string
	EnableUpnp       bool
	DatabaseFilePath string
}

func CreateConfigFile(outputFile string, settings SettingsShare) error {
	outputFolder := filepath.Dir(outputFile)
	os.MkdirAll(outputFolder, 0700)

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	err = e.Encode(settings)
	if err != nil {
		f.Close()
		return err
	}

	f.WriteString(buf.String())
	f.Close()

	return nil
}

func InitSettings(configFile string, settings *SettingsShare) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Info("New config file: %s\n", configFile)
		*settings = NewSettings()

		CreateConfigFile(configFile, *settings)
	} else {
		log.Info("Loading config file: %s\n", configFile)
		if _, err := toml.DecodeFile(configFile, &settings); err != nil {
			return err
		}
	}

	log.Info("Current config: %v\n", settings)
	return nil
}

func NewSettings() SettingsShare {
	return SettingsShare{Daemon: Daemon{Port: 7890, Host: "localhost", EnableUpnp: false, DatabaseFilePath: ConfigFileSQLITE()}}
}

func ConfigFolder() string {
	return UserHomeDir() + string(os.PathSeparator) + ".share"
}

func ConfigFile() string {
	return ConfigFolder() + string(os.PathSeparator) + "config.toml"
}

func ConfigFileSQLITE() string {
	return ConfigFolder() + string(os.PathSeparator) + "db.sqlite3"
}

func ConfigServerEndPoint(settings SettingsShare) string {
	return settings.Daemon.Host + ":" + strconv.Itoa(settings.Daemon.Port)
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			return os.Getenv("USERPROFILE")
		}

		return home
	}

	return os.Getenv("HOME")
}
