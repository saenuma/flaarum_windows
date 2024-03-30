package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/getlantern/systray"
)

//go:embed Logo.ico
var logoBytes []byte

//go:embed "artifacts"
var artifactsDir embed.FS
var StoreCancelFn context.CancelFunc

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	homeDir, _ := os.UserHomeDir()
	flaarumExecsPath := filepath.Join(homeDir, ".flaar312")
	if !DoesPathExists(flaarumExecsPath) {
		os.MkdirAll(flaarumExecsPath, 0777)

		dirFIs, err := artifactsDir.ReadDir("artifacts")
		if err != nil {
			log.Println(err)
			return
		}

		for _, dirFI := range dirFIs {
			dataOfFile, err := artifactsDir.ReadFile("artifacts/" + dirFI.Name())
			if err != nil {
				log.Println(err)
				return
			}

			os.WriteFile(filepath.Join(flaarumExecsPath, dirFI.Name()), dataOfFile, 0777)
		}

	} else {
		onSysVersionRaw, _ := os.ReadFile(filepath.Join(flaarumExecsPath, "version.txt"))

		versionRaw, _ := artifactsDir.ReadFile("artifacts/version.txt")
		if strings.TrimSpace(string(versionRaw)) != strings.TrimSpace(string(onSysVersionRaw)) {
			dirFIs, err := artifactsDir.ReadDir("artifacts")
			if err != nil {
				log.Println(err)
				return
			}

			for _, dirFI := range dirFIs {
				dataOfFile, err := artifactsDir.ReadFile("artifacts/" + dirFI.Name())
				if err != nil {
					log.Println(err)
					return
				}

				os.WriteFile(filepath.Join(flaarumExecsPath, dirFI.Name()), dataOfFile, 0777)
			}
		}
	}

	go func() {
		flStoreCtx, cancelFn := context.WithCancel(context.Background())
		StoreCancelFn = cancelFn
		absFlStorePath := filepath.Join(flaarumExecsPath, "flstore.exe")
		cmd := exec.CommandContext(flStoreCtx, absFlStorePath)
		cmd.Run()
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

	openFlaarumExecsFolder := systray.AddMenuItem("Open Flaarum Execs folder", "Open Flaarum Execs folder")
	go func() {
		<-openFlaarumExecsFolder.ClickedCh
		exec.Command("cmd", "/C", "start", flaarumExecsPath).Run()
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
	StoreCancelFn()
}

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
