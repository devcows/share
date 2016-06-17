package lib

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/BurntSushi/toml"
)

type SettingsShare struct {
	Daemon Daemon
}

type Daemon struct {
	Port int
	Host string
}

func InitSettings(settings *SettingsShare, portApi int) error {
	if _, err := os.Stat(ConfigFile()); os.IsNotExist(err) {
		fmt.Printf("New config file: %s\n", ConfigFile())
		os.MkdirAll(ConfigFolder(), 0700)
		settings = NewSettings()

		f, err := os.Create(ConfigFile())
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		e := toml.NewEncoder(&buf)
		err = e.Encode(settings)
		if err != nil {
			return err
		}

		f.WriteString(buf.String())
		f.Close()
	} else {
		fmt.Printf("Loading config file: %s\n", ConfigFile())
		if _, err := toml.DecodeFile(ConfigFile(), &settings); err != nil {
			return err
		}
	}

	settings.Daemon.Port = portApi
	return nil
}

func NewSettings() *SettingsShare {
	return &SettingsShare{Daemon: Daemon{Port: 7890, Host: "localhost"}}
}

func ConfigFolder() string {
	return UserHomeDir() + string(os.PathSeparator) + ".share"
}

func ConfigFile() string {
	return ConfigFolder() + string(os.PathSeparator) + "config.toml"
}

func ConfigFileSQLITE() string {
	return ConfigFolder() + string(os.PathSeparator) + "database.sqlite3"
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
