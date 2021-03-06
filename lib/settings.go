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
	ShareDaemon      ShareDaemon
	FileServerDaemon FileServerDaemon
	Mode             string
}

type FileServerDaemon struct {
	Port int
	Host string
}

type ShareDaemon struct {
	Port               int
	Host               string
	EnableUpnp         bool
	DatabaseFilePath   string
	CompressedFilePath string
}

const DIR_MASK = 0770

func CreateConfigFile(outputFile string, settings SettingsShare) error {
	outputFolder := filepath.Dir(outputFile)
	os.MkdirAll(outputFolder, DIR_MASK)

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	if err = e.Encode(settings); err != nil {
		return err
	}

	f.WriteString(buf.String())
	return nil
}

func InitSettings(configFile string, settings *SettingsShare) error {
	*settings = NewSettings()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.WithFields(log.Fields{"file": configFile}).Info("New config file.")

		if err := CreateConfigFile(configFile, *settings); err != nil {
			return err
		}
	} else {
		log.WithFields(log.Fields{"file": configFile}).Info("Loading config file.")

		if _, err := toml.DecodeFile(configFile, &settings); err != nil {
			return err
		}
	}

	os.MkdirAll(CompressedFilePath(), DIR_MASK)

	log.WithFields(log.Fields{"settings": settings}).Info("Current config.")
	return nil
}

func NewSettings() SettingsShare {
	settings := SettingsShare{Mode: "release"}

	settings.ShareDaemon = ShareDaemon{Port: 7890, Host: "localhost", EnableUpnp: false, DatabaseFilePath: ConfigFileSQLITE()}
	settings.FileServerDaemon = FileServerDaemon{Port: 7891, Host: ""}

	return settings
}

func ConfigFolder() string {
	return UserHomeDir(runtime.GOOS) + string(os.PathSeparator) + ".share"
}

func ConfigFile() string {
	return ConfigFolder() + string(os.PathSeparator) + "config.toml"
}

func CompressedFilePath() string {
	return ConfigFolder() + string(os.PathSeparator) + "files"
}

func ConfigFileSQLITE() string {
	return ConfigFolder() + string(os.PathSeparator) + "database.db"
}

func ConfigServerEndPoint(settings SettingsShare) string {
	return settings.ShareDaemon.Host + ":" + strconv.Itoa(settings.ShareDaemon.Port)
}

func ConfigFileServerEndPoint(settings SettingsShare) string {
	return settings.FileServerDaemon.Host + ":" + strconv.Itoa(settings.FileServerDaemon.Port)
}

func UserHomeDir(runtime_goos string) string {
	if runtime_goos == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			return os.Getenv("USERPROFILE")
		}

		return home
	}

	return os.Getenv("HOME")
}
