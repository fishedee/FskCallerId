package main

import (
	"fmt"
	. "fskcallerid/util"
	. "github.com/fishedee/app/workgroup"
	. "github.com/fishedee/language"
	"os"
	"strconv"
	"time"
)

type Server struct {
	log        Log
	serial     *Serial
	subscriber *Subscriber
	contact    *Contact
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
	data := map[string]interface{}{
		"type":   "call",
		"time":   callTime,
		"phone":  phone,
		"name":   contactInfo.Name,
		"image":  "",
		"remark": fmt.Sprintf("[%v]%v", contactInfo.Group, contactInfo.Remark),
	}
	//打印出来
	this.log.Debug("call in one!%v", data)
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
	//打印出来
	this.log.Debug("hang in one!")
}

func (this *Server) Run() error {
	this.contact.Init()
	this.serial.SetOnCall(this.callListener)
	this.serial.SetOnHang(this.hangListener)
	this.subscriber.SetContactUpdateListener(this.contactUpdateListener)

	this.log.Debug("server is running...")

	return nil
}

func (this *Server) Close() {

}

func init() {
	MustRegisterIoc(NewServer)
}

func main() {
	defer CatchCrash(func(e Exception) {
		fmt.Printf("init fail!%v\n", e.Error())
	})

	err := os.Chdir("..")
	if err != nil {
		panic(err)
	}
	MustInvokeIoc(func(log Log, queue Queue, server *Server, serial *Serial, subscriber *Subscriber) {
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
