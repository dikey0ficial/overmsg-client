package main

import (
	"gioui.org/x/pref/theme"
	"github.com/BurntSushi/toml"
	"github.com/sqweek/dialog"
	"io/ioutil"
	"os"
)

var conf = struct {
	Name       string   `toml:"name"`
	Token      string   `toml:"token"`
	IsDark     bool     `toml:"is_dark"`
	ServerURLs []string `toml:"server_urls"`
}{}

func initConfig() {
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
	if (conf.Name == "" && conf.Token != "") || (conf.Name != "" && conf.Token == "") {
		conf.Name, conf.Token = "", ""
		saveConf() // no checking error because yes))))
	}
	if len(conf.ServerURLs) == 0 {
		conf.ServerURLs = []string{
			// while i haven't deployed server, there will be only localhost
			"localhost",
		}
		saveConf() // the same as in last if
	}
}

func saveConf() error {
	f, err := os.Create("config.toml")
	if err != nil {
		return err
	}
	return toml.NewEncoder(f).Encode(conf)
}
