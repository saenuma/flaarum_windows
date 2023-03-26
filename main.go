package main

import (
	_ "embed"
	"fmt"

	"github.com/getlantern/systray"
)

//go:embed Logo.ico
var logoBytes []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	fmt.Println(len(logoBytes))
	systray.SetIcon(logoBytes)
	systray.SetTitle("Flaarum")
	systray.SetTooltip("Flaarum: a more comfortable database")
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
