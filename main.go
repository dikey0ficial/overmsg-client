package main

import (
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
	errl = log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func work(ui *UI) {
	fsize := [2]unit.Value{unit.Dp(512), unit.Dp(512)}
	options := []app.Option{
		app.Title("OVERMSg"),
		app.Size(fsize[0], fsize[1]),
		app.MinSize(fsize[0], fsize[1]),
	}
	w := app.NewWindow(options...)
	if err := ui.Run(w); err != nil {
		dialog.Message("Error: %v", err).Title("Error!!1").Error()
		errl.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	ui := NewUI()
	go work(ui)
	app.Main()
}
