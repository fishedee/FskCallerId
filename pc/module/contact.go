package module

import (
	"errors"
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/app/queue"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/util"
	"io/ioutil"
	"os"
	"time"
)

type ContactInfo map[string]ContactSingleInfo

type ContactSingleInfo struct {
	Name   string
	Remark string
	Group  string
}

type VersionInfo struct {
	ContactVersionId int
	CreateTime       time.Time
	ModifyTIme       time.Time
}

type DataInfo struct {
	Contact ContactInfo
	Version VersionInfo
}

type Contact struct {
	log      Log
	queue    Queue
	config   ContactConfig
	data     *DataInfo
	ajaxPool *AjaxPool
}

type ContactConfig struct {
	Server string `config:"server"`
	File   string `config:"file"`
}

func NewContact(log Log, queue Queue, config ContactConfig) (*Contact, error) {
	contact := &Contact{
		log:      log,
		queue:    queue,
		config:   config,
		ajaxPool: NewAjaxPool(&AjaxPoolOption{}),
	}
	queue.MustConsume("/contact/checkVersion", "contactAo", 1, contact.checkVersion)
	return contact, nil
}

func (this *Contact) Init() {
	data, err := ioutil.ReadFile(this.config.File)
	if err != nil {
		panic(err)
	}
	err = DecodeJson([]byte(data), &this.data)
	if err != nil {
		panic("invalid database format " + err.Error())
	}
	this.CheckVersion(-1)
}

func (this *Contact) GetContactByPhone(phone string) ContactSingleInfo {
	contactSingleInfo, isExist := this.data.Contact[phone]
	if isExist == true {
		return contactSingleInfo
	}
	return ContactSingleInfo{
		Name:   phone,
		Remark: "",
		Group:  "未分组",
	}
}

func (this *Contact) GetVersion() VersionInfo {
	return this.data.Version
}

func (this *Contact) CheckVersion(versionId int) {
	this.queue.MustProduce("/contact/checkVersion", versionId)
}

func (this *Contact) checkVersion(versionId int) {
	curVersionId := this.data.Version.ContactVersionId
	if versionId > 0 && curVersionId >= versionId {
		return
	}

	//拉取最新的数据
	var dataByte []byte
	err := this.ajaxPool.Get(&Ajax{
		Url: this.config.Server + "/contact/getall",
		Header: map[string]string{
			"Accept-Encoding": "gzip",
		},
		ResponseData: &dataByte,
	})
	if err != nil {
		panic(err)
	}
	var data struct {
		Code int
		Msg  string
		Data DataInfo
	}
	err = DecodeJson(dataByte, &data)
	if err != nil {
		panic(err)
	}
	if data.Code != 0 {
		panic(data.Msg)
	}
	//更新内存数据
	this.data = &data.Data
	//更新数据库
	dataByte, err = EncodeJson(data.Data)
	if err != nil {
		panic(err)
	}
	tempFile := this.config.File + "_temp"
	err = ioutil.WriteFile(tempFile, dataByte, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.Rename(tempFile, this.config.File)
	if err != nil {
		panic(err)
	}
}

func (this *Contact) SendCallLog(phone string, callTime time.Time) error {
	var dataByte []byte
	err := this.ajaxPool.Post(&Ajax{
		Url: this.config.Server + "/calllog/add",
		Data: map[string]string{
			"phone":    phone,
			"callTime": callTime.Format("2006-01-02 15:04:05"),
		},
		ResponseData: &dataByte,
	})
	if err != nil {
		return err
	}
	var data struct {
		Code int
		Msg  string
		Data VersionInfo
	}
	err = DecodeJson(dataByte, &data)
	if err != nil {
		return err
	}
	if data.Code != 0 {
		return errors.New("/calllog/add fail:" + data.Msg)
	}
	this.CheckVersion(data.Data.ContactVersionId)
	return nil
}
