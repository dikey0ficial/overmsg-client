package main

import (
	"gioui.org/x/pref/theme"
	"github.com/BurntSushi/toml"
	"github.com/sqweek/dialog"
	"io/ioutil"
	"os"
)

var conf = struct {
	Name   string `toml:"name"`
	Token  string `toml:"token"`
	IsDark bool   `toml:"is_dark"`
}{}

func init() {
	f, err := os.Open("config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create("config.toml")
			conf.IsDark, _ = theme.IsDarkMode()
			if err != nil {
				dialog.Message("Error: %v", err).Title("Error!!1").Error()
				errl.Println(err)
				os.Exit(1)
			}
		} else {
			dialog.Message("Error: %v", err).Title("Error!!1").Error()
			errl.Println(err)
			os.Exit(1)
		}
	}
	defer f.Close()
	dat, err := ioutil.ReadAll(f)
	if err != nil {
		dialog.Message("Error: %v", err).Title("Error!!1").Error()
		errl.Println(err)
		os.Exit(1)
	}
	if err := toml.Unmarshal(dat, &conf); err != nil {
		dialog.Message("Error: %v", err).Title("Error!!1").Error()
		errl.Println(err)
		os.Exit(1)
	}
}

func saveConf() error {
	f, err := os.Create("config.toml")
	if err != nil {
		return err
	}
	return toml.NewEncoder(f).Encode(conf)
}
