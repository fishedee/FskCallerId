package main

import (
	"fmt"
	. "fskcallerid/util"
	. "github.com/fishedee/app/workgroup"
	. "github.com/fishedee/language"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Server struct {
	log          Log
	serial       *Serial
	subscriber   *Subscriber
	contact      *Contact
	tooltip      *NotifyIcon
	callerWindow *CallerWindow
}

func NewServer(log Log, serial *Serial, contact *Contact, subscriber *Subscriber) *Server {
	server := &Server{}
	server.log = log
	server.serial = serial
	server.contact = contact
	server.subscriber = subscriber
	return server
}

func (this *Server) callListener(callTime time.Time, phone string) {
	//记录通讯记录
	go func() {
		err := this.contact.SendCallLog(phone, callTime)
		if err != nil {
			this.log.Critical("SendCallLog error %v", err.Error())
		}
	}()

	//读取电话对应信息
	contactInfo := this.contact.GetContactByPhone(phone)
	/*
		data := map[string]interface{}{
			"type":   "call",
			"time":   callTime,
			"phone":  phone,
			"name":   contactInfo.Name,
			"image":  "",
			"remark": fmt.Sprintf("[%v]%v", contactInfo.Group, contactInfo.Remark),
		}
	*/

	//打开窗口
	this.tooltip.Synchronize(func() {
		//先关掉原来的窗口
		if this.callerWindow != nil {
			this.callerWindow.Dispose()
			this.callerWindow = nil
		}
		//再创建新窗口
		this.callerWindow = NewCallerWindow()
		this.callerWindow.SetCaller(callTime.Format("15:04"), phone, contactInfo.Name, "["+contactInfo.Group+"]"+contactInfo.Remark)
		this.callerWindow.Show()
	})
}

func (this *Server) contactUpdateListener(version string) {
	//检查当前的版本
	versionId, err := strconv.Atoi(version)
	if err != nil {
		panic(err)
	}
	this.contact.CheckVersion(versionId)
}

func (this *Server) hangListener() {
	//关闭窗口
	this.tooltip.Synchronize(func() {
		if this.callerWindow != nil {
			this.callerWindow.Dispose()
			this.callerWindow = nil
		}
	})
}

func (this *Server) Run() error {
	this.contact.Init()
	this.serial.SetOnCall(this.callListener)
	this.serial.SetOnHang(this.hangListener)
	this.subscriber.SetContactUpdateListener(this.contactUpdateListener)

	this.callerWindow = nil
	this.tooltip = NewNotifyIcon()

	this.tooltip.AddAction("退出", func() {
		os.Exit(0)
	})

	this.tooltip.AddAction("显示版本", func() {
		version := this.contact.GetVersion()
		title := "电话目录版本"
		message := fmt.Sprintf("版本号:%v,时间:%v", version.ContactVersionId, version.CreateTime.Local().Format("2006-01-02 15:04:05"))
		this.tooltip.MessageBox(title, message)
	})

	this.log.Debug("server is running...")

	this.tooltip.Run()

	return nil
}

func (this *Server) Close() {

}

func init() {
	MustRegisterIoc(NewServer)
	MustRegisterIoc(func() ConfigFile {
		return "./config.ini"
	})
}

func main() {
	//FIXME fmt在Windows下打印不出来
	defer CatchCrash(func(e Exception) {
		fmt.Printf("init fail!%v\n", e.Error())
	})

	//切换当前目录
	currentDir := filepath.Dir(os.Args[0])
	err := os.Chdir(currentDir)
	if err != nil {
		panic(err)
	}

	MustInvokeIoc(func(log Log, queue Queue, server *Server, serial *Serial, subscriber *Subscriber) {
		defer CatchCrash(func(e Exception) {
			log.Critical("%v", e.Error())
		})

		workgroup, err := NewWorkGroup(log, WorkGroupConfig{})
		if err != nil {
			panic(err)
		}
		workgroup.Add(serial)
		workgroup.Add(subscriber)
		workgroup.Add(queue)
		workgroup.Add(server)
		err = workgroup.Run()
		if err != nil {
			panic(err)
		}
	})
}
