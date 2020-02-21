package main

import (
	"github.com/lxn/walk"
)

type NotifyIcon struct {
	walk.NotifyIcon
}

func NewNotifyIcon(mw walk.Form) *NotifyIcon {
	//创建notifyIcon
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		panic(err)
	}

	icon, err := walk.Resources.Icon("3")
	if err != nil {
		panic(err)
	}
	err = ni.SetIcon(icon)
	if err != nil {
		panic(err)
	}

	if err := ni.SetToolTip("来电提醒系统已经运行中"); err != nil {
		panic(err)
	}

	//右键菜单，退出按钮
	exitAction := walk.NewAction()
	if err := exitAction.SetText("退出"); err != nil {
		panic(err)
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		panic(err)
	}
	if err := ni.SetVisible(true); err != nil {
		panic(err)
	}
	/*
		if err := ni.ShowInfo("英豪彩瓦厂", "来电提醒系统已经启动..."); err != nil {
			panic(err)
		}
	*/

	return &NotifyIcon{
		NotifyIcon: *ni,
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
