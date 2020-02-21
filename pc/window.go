package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type CallerWindow struct {
	window *walk.MainWindow
}

func NewCallerWindow() *CallerWindow {
	var window *walk.MainWindow
	err := MainWindow{
		AssignTo: &window,
		Name:     "fisdh",
		Title:    "fishedee来电显示",
		Size:     Size{300, 150},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{},
	}.Create()
	if err != nil {
		panic(err)
	}

	return &CallerWindow{
		window: window,
	}
}

func (this *CallerWindow) GetWindow() *walk.MainWindow {
	return this.window
}

func (this *CallerWindow) Run() {
	code := this.window.Run()
	fmt.Println(code)
}

func (this *CallerWindow) Show() {
	hWnd := this.window.Handle()
	win.SetWindowLong(hWnd, win.GWL_STYLE, win.WS_OVERLAPPED|win.WS_CAPTION|win.WS_SYSMENU)

	var mi win.MONITORINFO
	moniter := win.MonitorFromWindow(hWnd, win.MONITOR_DEFAULTTOPRIMARY)
	win.GetMonitorInfo(moniter, &mi)
	moniterSize := mi.RcMonitor
	win.SetWindowPos(hWnd, win.HWND_TOP, moniterSize.Right-300, moniterSize.Bottom-150, 300, 150, win.SWP_FRAMECHANGED|win.SWP_NOOWNERZORDER)

	this.window.Show()
}

func (this *CallerWindow) Hide() {
	this.window.Hide()
}
