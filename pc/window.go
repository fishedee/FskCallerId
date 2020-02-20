package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func main() {
	//walk.MainWindow是MainWindow的实例
	//declarative.MainWindow是MainWindow的声明式创建
	var mainWindow *walk.MainWindow
	err := MainWindow{
		AssignTo: &mainWindow,
		Title:    "fishedee来电显示",
		Size:     Size{300, 150},
	}.Create()
	if err != nil {
		panic(err)
	}
	mainWindow.Run()

	//消失与显示
	//mainWindow.Show()
	//mainWindow.Hide()

	//设置标题栏样式，可以重叠，有标题栏，有系统菜单，没有边框拉伸，最大化和最小化按钮
	hWnd := mainWindow.Handle()
	win.SetWindowLong(hWnd, win.GWL_STYLE, WS_OVERLAPPED|WS_CAPTION|WS_SYSMENU)

	//将窗口显示在右下角
	var mi win.MONITORINFO
	moniter := win.MonitorFromWindow(hWnd, win.MONITOR_DEFAULTTOPRIMARY)
	win.GetMonitorInfo(moniter, &mi)
	moniterSize := mi.RcMonitor
	win.SetWindowPos(hWnd, win.HWND_TOP, moniterSize.Right-300, moniterSize.Bottom-150, 300, 150, win.SWP_FRAMECHANGED|win.SWP_NOOWNERZORDER)

	//跨线程投递消息到窗口
	lv.logChan <- value
	win.PostMessage(hWnd, win.WM_USER, 0, 0)

}

type MyMainWindow struct {
	walk.MainWindow
}

func (this *MyMainWindow) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_USER:
		select {
		case value := <-lv.logChan:
			//dosomething you should
		default:
			return 0
		}
	}

	return this.MainWindow.WndProc(hwnd, msg, wParam, lParam)
}
