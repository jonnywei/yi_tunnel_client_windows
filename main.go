package main

import (
	"fmt"
	"github.com/getlantern/systray"
	yitunnel "github.com/jonnywei/yi_tunnel/client"
	"github.com/jonnywei/yi_tunnel_client_windows/icon"
	"io"
	"log"
	"os"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("YiTunnelWindows")
	systray.SetTooltip("YiTunnelWindows Free You")
	config := yitunnel.LoadConfigFile("./config.json")
	mGroup := systray.AddMenuItemCheckbox("启动Yi", "启动本地节点", false)
	mExitYi := systray.AddMenuItemCheckbox("关闭Yi", "关闭本地节点", true)

	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("退出程序", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
		os.Exit(1)
	}()
	go func() {
		for {
			select {
			case <-mGroup.ClickedCh:
				fmt.Println("Requesting 启动")
				yitunnel.RunClient(config)
				fmt.Println("Finished 启动本地节点")
				mGroup.Check()
				mExitYi.Uncheck()
				systray.SetTemplateIcon(icon.ConnectedData, icon.ConnectedData)
			case <-mExitYi.ClickedCh:
				fmt.Println("Requesting 关闭")
				yitunnel.Close()
				fmt.Println("Finished 关闭")
				mExitYi.Check()
				mGroup.Uncheck()
				systray.SetTemplateIcon(icon.Data, icon.Data)
			}
		}
	}()
}

func onExit() {
	fmt.Println("Exist Program ")
	yitunnel.Close()
	fmt.Println("Exist Program Finish")
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	s := "you"
	fmt.Printf("Hello and welcome, %s!.You are Free ...\n", s)
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	defer logFile.Close()
	systray.Run(onReady, onExit)
}
