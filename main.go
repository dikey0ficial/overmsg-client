package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/sqweek/dialog"
	"log"
	"os"
)

var (
	debl, errl *log.Logger
)

func init() {
	debl = log.New(os.Stdout, "[DEBUG]\t", log.Ldate|log.Ltime|log.Lshortfile)
	errlf, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		dialog.Message("Error openning errors log))) (error: %v)", err).Title("Error!!1").Error()
		os.Exit(1)
	}
	errl = log.New(errlf, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	initConfig()
	initAPI()
}

func work(ui *UI) {
	defer func() {
		if e := recover(); e != nil {
			errl.Println(e)
			dialog.Message("Fatal error :(((((((((((").Title("ERROR!!!!!!!!!!!!!!!").Error()
			os.Exit(1)
		}
	}()
	fsize := [2]unit.Value{unit.Dp(768), unit.Dp(512)}
	options := []app.Option{
		app.Title("OVERMSg"),
		app.Size(fsize[0], fsize[1]),
		app.MinSize(fsize[0], fsize[1]),
	}
	w := app.NewWindow(options...)
	if err := ui.Run(w); err != nil {
		if err == errSAW {
			w = nil
		}
		errl.Println(err)
		dialog.Message("Error: %v", err).Title("Error!!1").Error()
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	ui := NewUI()
	go work(ui)
	app.Main()
}
