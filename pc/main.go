package main

import (
	"fmt"
	. "github.com/fishedee/language"
)

func main() {
	defer func() {
		recover()
		for {
		}
	}()
	defer CatchCrash(func(e Exception) {
		fmt.Println(e)
	})
	window := NewCallerWindow()
	window.Show()

	tooltip := NewNotifyIcon()

	defer tooltip.Dispose()

	tooltip.AddAction("显示电话", func() {
		//window.SetCaller()
	})

	tooltip.AddAction("显示窗口", func() {
		window.SetRing()
	})

	tooltip.Run()
}
