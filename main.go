package main

import (
	_ "embed"
	"fmt"
	"os/exec"

	"github.com/getlantern/systray"
)

//go:embed Logo.ico
var logoBytes []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(logoBytes)
	systray.SetTitle("Flaarum")
	systray.SetTooltip("Flaarum: a more comfortable database")

	flaarumTuts := systray.AddMenuItem("Flaarum tutorials", "Launch flaarum tutorials")
	go func() {
		<-flaarumTuts.ClickedCh
		exec.Command("cmd", "/C", "start", "https://sae.ng/flaarumtuts").Run()
	}()

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

}

func onExit() {
	// clean up here
}
