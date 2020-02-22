package main

import (
	"github.com/lxn/walk"
)

func main() {

	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		panic(err)
	}

	window := NewCallerWindow()
	window.SetCaller()
	window.Show()

	tooltip := NewNotifyIcon(mainWindow)

	defer tooltip.Dispose()

	tooltip.AddAction("显示窗口", func() {

	})

	mainWindow.Run()
}
