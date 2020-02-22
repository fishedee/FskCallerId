package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"unsafe"
)

var (
	HEIGHT int32 = 300
	WIDTH  int32 = 500
)

type CallerWindow struct {
	window    *walk.MainWindow
	container *walk.Composite
}

func NewCallerWindow() *CallerWindow {

	var window *walk.MainWindow
	var container *walk.Composite
	err := MainWindow{
		AssignTo:   &window,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Name:       "",
		Title:      "来电显示提醒",
		Size:       Size{int(WIDTH), int(HEIGHT)},
		Layout:     VBox{Margins: Margins{14, 14, 14, 14}, Spacing: 20},
		Children: []Widget{
			Composite{
				MinSize:    Size{0, 20},
				Background: SolidColorBrush{Color: walk.RGB(105, 105, 105)},
				Layout:     HBox{},
			},
			Composite{
				AssignTo:   &container,
				Background: SolidColorBrush{Color: walk.RGB(255, 152, 29)},
				Layout:     HBox{MarginsZero: true, Spacing: 20},
				Children: []Widget{
					ImageView{
						Background: SolidColorBrush{Color: walk.RGB(0, 64, 96)},
						Image:      "head.jpg",
						Mode:       ImageViewModeShrink,
						MaxSize:    Size{120, 120},
					},
					Composite{
						Background: SolidColorBrush{Color: walk.RGB(0, 130, 135)},
						Alignment:  AlignHNearVNear,
						Layout:     VBox{MarginsZero: true, SpacingZero: true},
						Children: []Widget{
							Composite{
								Background: SolidColorBrush{Color: walk.RGB(255, 0, 0)},
								Layout:     HBox{MarginsZero: true, SpacingZero: true},
								Children: []Widget{
									TextLabel{Text: "123", Font: Font{PointSize: 16}},
									VSpacer{},
									TextLabel{Text: "123", Font: Font{PointSize: 16}},
								},
							},
							TextLabel{Text: "123", Font: Font{PointSize: 16}},
							VSpacer{Size: 40},
							TextLabel{Text: "456", Font: Font{PointSize: 16}},
						},
					},
				},
			},
			Composite{
				MinSize: Size{0, 20},
				Layout:  HBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					Composite{
						Background:    SolidColorBrush{Color: walk.RGB(0, 64, 96)},
						StretchFactor: 0,
						Layout:        HBox{},
					},
					Composite{
						Background:    SolidColorBrush{Color: walk.RGB(0, 130, 135)},
						StretchFactor: 0,
						Layout:        HBox{},
					},
					Composite{
						Background:    SolidColorBrush{Color: walk.RGB(255, 152, 29)},
						StretchFactor: 0,
						Layout:        HBox{},
					},
					Composite{
						Background:    SolidColorBrush{Color: walk.RGB(224, 100, 135)},
						StretchFactor: 0,
						Layout:        HBox{},
					},
					Composite{
						Background:    SolidColorBrush{Color: walk.RGB(120, 186, 0)},
						StretchFactor: 0,
						Layout:        HBox{},
					},
				},
			},
		},
		Visible: false,
	}.Create()
	if err != nil {
		panic(err)
	}

	//设置图标
	icon, err := walk.Resources.Icon("3")
	if err != nil {
		panic(err)
	}
	err = window.SetIcon(icon)
	if err != nil {
		panic(err)
	}

	//设置窗口样式
	hWnd := window.Handle()
	win.SetWindowLong(hWnd, win.GWL_STYLE, win.WS_OVERLAPPED|win.WS_CAPTION|win.WS_SYSMENU)
	return &CallerWindow{
		window:    window,
		container: container,
	}
}

func (this *CallerWindow) GetWindow() *walk.MainWindow {
	return this.window
}

func (this *CallerWindow) SetCaller() {
	this.container.Children().At(1).Dispose()
	builder := NewBuilder(this.container)

	Composite{
		Background:    SolidColorBrush{Color: walk.RGB(0, 130, 135)},
		Layout:        VBox{},
		StretchFactor: 20,
		Children: []Widget{
			Label{Text: "ring..."},
			Label{Text: "ring..."},
			Label{Text: "ring..."},
		},
	}.Create(builder)
}

func (this *CallerWindow) SetRing() {
	this.container.Children().At(1).Dispose()
	builder := NewBuilder(this.container)

	Composite{
		Layout:        VBox{},
		StretchFactor: 20,
		Children:      []Widget{
			//Label{Text: "ring..."},
		},
	}.Create(builder)
}

func (this *CallerWindow) Show() {
	hWnd := this.window.Handle()

	var rect win.RECT

	win.SystemParametersInfo(48, 0, unsafe.Pointer(&rect), 0)

	win.SetWindowPos(hWnd, win.HWND_TOP, rect.Right-WIDTH, rect.Bottom-HEIGHT, WIDTH, HEIGHT, win.SWP_FRAMECHANGED|win.SWP_NOOWNERZORDER)

	this.window.Show()
}

func (this *CallerWindow) Hide() {
	this.window.Hide()
}

func (this *CallerWindow) Dispose() {
	this.window.Dispose()
}
