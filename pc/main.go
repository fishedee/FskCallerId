package main

import (
	"github.com/lxn/walk"
	"log"
)

func GuiInit() {
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}
	//托盘图标文件
	icon, err := walk.Resources.Icon("./caller.ico")
	if err != nil {
		log.Fatal(err)
	}
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("鼠标在icon上悬浮的信息."); err != nil {
		log.Fatal(err)
	}
	ni.MouseDown().Attach(func(x, yint, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		if err := ni.ShowCustom("Walk 任务栏通知标题", "walk 任务栏通知内容"); err != nil {
			og.Fatal(err)
		}
	})
	exitAction := walk.NewAction()
	if err := exitAction.SetText("右键icon的菜单按钮"); err != nil {
		log.Fatal(err)
	}
	//Exit 实现的功能
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}
	if err := ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again."); err != nil {
		log.Fatal(err)
	}
	mw.Run()
}
