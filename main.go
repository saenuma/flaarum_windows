package main

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/getlantern/systray"
)

//go:embed Logo.ico
var logoBytes []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go func() {
		absFlStorePath, _ := filepath.Abs("flstore.exe")
		exec.Command(absFlStorePath).Run()
	}()

	systray.SetIcon(logoBytes)
	systray.SetTitle("Flaarum")
	systray.SetTooltip("Flaarum: a more comfortable database")

	flaarumTuts := systray.AddMenuItem("Flaarum tutorials", "Launch flaarum tutorials")
	go func() {
		<-flaarumTuts.ClickedCh
		exec.Command("cmd", "/C", "start", "https://sae.ng/flaarumtuts").Run()
	}()
	systray.AddSeparator()

	openFlaarumFolder := systray.AddMenuItem("Open Flaarum folder", "Open Flaarum folder")
	go func() {
		<-openFlaarumFolder.ClickedCh
		hd, _ := os.UserHomeDir()
		flaarumFolder := filepath.Join(hd, "Flaarum")
		os.MkdirAll(flaarumFolder, 0777)
		exec.Command("cmd", "/C", "start", flaarumFolder).Run()
	}()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

}

func onExit() {
	// clean up here
}
