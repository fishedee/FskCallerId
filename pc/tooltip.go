package main

import (
	"github.com/lxn/walk"
)

type NotifyIcon struct {
	walk.NotifyIcon
	mainWindow *walk.MainWindow
}

func NewNotifyIcon() *NotifyIcon {
	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		panic(err)
	}

	//创建notifyIcon
	ni, err := walk.NewNotifyIcon(mainWindow)
	if err != nil {
		panic(err)
	}

	err = ni.SetIcon(callerIcon)
	if err != nil {
		panic(err)
	}

	if err := ni.SetToolTip("来电提醒系统已经运行中"); err != nil {
		panic(err)
	}

	if err := ni.SetVisible(true); err != nil {
		panic(err)
	}
	if err := ni.ShowInfo("英豪彩瓦厂", "来电提醒系统已经启动..."); err != nil {
		panic(err)
	}

	return &NotifyIcon{
		NotifyIcon: *ni,
		mainWindow: mainWindow,
	}
}

func (this *NotifyIcon) Dispose() {
	this.NotifyIcon.Dispose()
}

func (this *NotifyIcon) AddAction(name string, handler func()) {
	action := walk.NewAction()
	if err := action.SetText(name); err != nil {
		panic(err)
	}
	action.Triggered().Attach(handler)
	if err := this.NotifyIcon.ContextMenu().Actions().Add(action); err != nil {
		panic(err)
	}
}

func (this *NotifyIcon) Run() {
	this.mainWindow.Run()
}

func (this *NotifyIcon) MessageBox(title string, message string) {
	walk.MsgBox(this.mainWindow, title, message, walk.MsgBoxOK)
}

func (this *NotifyIcon) Synchronize(handler func()) {
	this.mainWindow.Synchronize(handler)
}
