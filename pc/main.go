package main

import (
	"fmt"
	"github.com/lxn/walk"
)

func main() {

	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		panic(err)
	}

	window := NewCallerWindow()

	tooltip := NewNotifyIcon(mainWindow)

	defer tooltip.Dispose()

	tooltip.AddAction("显示窗口", func() {
		window.Show()
	})

	fmt.Println("c1")

	mainWindow.Run()

	fmt.Println("c2")
}
