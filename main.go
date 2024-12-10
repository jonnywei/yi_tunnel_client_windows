package main

import (
	"fmt"
	"github.com/getlantern/systray"
	yitunnel "github.com/jonnywei/yi_tunnel"
	"github.com/jonnywei/yi_tunnel_client_windows/icon"
	"io"
	"log"
	"os"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("GolandProjects")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
	mGroup := systray.AddMenuItem("启动", "启动本地节点")
	go func() {

		<-mGroup.ClickedCh
		fmt.Println("Requesting 启动")
		yitunnel.RunClient()
		fmt.Println("Finished 启动本地节点")
	}()
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	s := "gopher"
	fmt.Println("Hello and welcome, %s!", s)

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	defer logFile.Close()
	systray.Run(onReady, func() {})
}
