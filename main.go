package main

import (
	"errors"
	"fmt"
	"github.com/getlantern/systray"
	yitunnel "github.com/jonnywei/yi_tunnel/client"
	"github.com/jonnywei/yi_tunnel_client_windows/icon"
	"golang.org/x/sys/windows/registry"
	"io"
	"log"
	"os"
)

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	programKey := "YiTunnelWindows"
	systray.SetTitle(programKey)
	systray.SetTooltip("YiTunnelWindows Free You")
	config := yitunnel.LoadConfigFile("./config.json")
	mGroup := systray.AddMenuItemCheckbox("启动Yi", "启动本地节点", false)
	mExitYi := systray.AddMenuItemCheckbox("关闭Yi", "关闭本地节点", true)
	systray.AddSeparator()
	mAutoStart := systray.AddMenuItemCheckbox("开机启动", "开机启动", false)
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
		auto, err := isAutoStartWithWindows(programKey)
		if err != nil {
			fmt.Println(err)
			mAutoStart.Uncheck()
			return
		}
		if auto {
			mAutoStart.Check()
		}
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
			case <-mAutoStart.ClickedCh:
				if mAutoStart.Checked() {
					fmt.Println("remove auto start ")
					_, err := removeAutoStartWithWindows(programKey)
					if err == nil {
						mAutoStart.Uncheck()
					}
				} else {
					fmt.Println("set Auto start ")
					_, err := setAutoStartWithWindows(programKey)
					if err == nil {
						mAutoStart.Check()
					}
				}
			}
		}
	}()
}

func onExit() {
	fmt.Println("Exist Program ")
	yitunnel.Close()
	fmt.Println("Exist Program Finish")
}

func isAutoStartWithWindows(programKey string) (bool, error) {

	exePath, err := os.Executable()
	if err != nil {
		log.Printf("os.Executable: %v", err)
		return false, err
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Printf("error registry.OpenKey: %v", err)
		return false, err
	}
	defer key.Close()
	_, _, err = key.GetStringsValue(programKey)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		log.Printf("already set key: %v", exePath)
	}
	return true, nil
}

func setAutoStartWithWindows(programKey string) (bool, error) {
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("os.Executable: %v", err)
		return false, err
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Printf("error registry.OpenKey: %v", err)
	}
	defer key.Close()
	_, _, err = key.GetStringsValue(programKey)
	if errors.Is(err, os.ErrNotExist) {
		err = key.SetStringValue(programKey, exePath)
		if err != nil {
			log.Printf("error key.SetStringsValue: %v", err)
			return false, err
		}
		log.Printf("set run: %v", exePath)
		return true, nil
	} else {
		log.Printf("already set key: %v", exePath)
	}
	return true, nil
}

func removeAutoStartWithWindows(programKey string) (bool, error) {

	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Printf("error registry.OpenKey: %v", err)
		return false, err
	}
	defer key.Close()
	err = key.DeleteValue(programKey)
	if err != nil {
		return false, err
	}
	return true, nil
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
