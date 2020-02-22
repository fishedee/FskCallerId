package module

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/language"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"strconv"
	"strings"
	"time"
)

type Serial struct {
	msgChan      chan string
	log          Log
	config       SerialConfig
	callListener SerialCallListener
	hangListener SerialHangListener
	hangTimeout  <-chan time.Time
}

type SerialCallListener func(time time.Time, phone string)

type SerialHangListener func()

type SerialConfig struct {
	PortName      string        `config:"portname"`
	BaudRate      int           `config:"baudrate"`
	RingAfterHang time.Duration `config:"ringafterhang"`
}

func NewSerial(log Log, config SerialConfig) (*Serial, error) {
	serial := &Serial{}
	serial.log = log
	serial.config = config
	serial.msgChan = make(chan string, 256)
	serial.hangTimeout = nil
	return serial, nil
}

func (this *Serial) SetOnCall(callListener SerialCallListener) {
	this.callListener = callListener
}

func (this *Serial) SetOnHang(hangListener SerialHangListener) {
	this.hangListener = hangListener
}

func (this *Serial) parseCallerInfo(msg string) (time.Time, string) {
	if len(msg) < 12 {
		panic("invalid msg" + msg)
	}
	msg = msg[4:]
	this.log.Informational("caller info: %v,%v", msg, len(msg))
	month, err := strconv.Atoi(msg[:2])
	if err != nil {
		panic(err)
	}
	day, err := strconv.Atoi(msg[2:4])
	if err != nil {
		panic(err)
	}
	hour, err := strconv.Atoi(msg[4:6])
	if err != nil {
		panic(err)
	}
	min, err := strconv.Atoi(msg[6:8])
	if err != nil {
		panic(err)
	}
	now := time.Now()
	callerTime := time.Date(now.Year(), time.Month(month), day, hour, min, 0, 0, now.Location())
	callerPhone := msg[8:]
	return callerTime, callerPhone
}

func (this *Serial) fireSingleListener() {
	defer CatchCrash(func(e Exception) {
		this.log.Critical("serial crash fire %v", e.Error())
	})

	select {
	case msg := <-this.msgChan:
		this.log.Debug("receive serial info: %v,%v", msg, len(msg))
		if msg == "ATLINKOK" {
			//刚连接成功的消息
			this.log.Debug("serial link connect ok!")
		} else if msg == "ATRING" {
			//报出ring信息
			if this.hangTimeout != nil {
				this.hangTimeout = time.After(this.config.RingAfterHang)
			}
		} else if strings.Index(msg, "ATF") == 0 {
			//收到号码信息
			callerTime, callerPhone := this.parseCallerInfo(msg)
			this.log.Debug("receive serial caller info:%v,%v", callerTime, callerPhone)
			if this.hangTimeout != nil {
				if this.hangListener != nil {
					this.hangListener()
				}
				this.hangTimeout = nil
			}
			this.hangTimeout = time.After(this.config.RingAfterHang)
			if this.callListener != nil {
				this.callListener(callerTime, callerPhone)
			}
		} else {
			panic("invalid msg " + msg)
		}
	case <-this.hangTimeout:
		//超时，看作已经挂断了
		if this.hangListener != nil {
			this.hangListener()
		}
		this.hangTimeout = nil
	}
}

func (this *Serial) fireListener() {
	for {
		this.fireSingleListener()
	}
}

func (this *Serial) singleRun() {
	options := serial.OpenOptions{
		PortName:              this.config.PortName,
		BaudRate:              uint(this.config.BaudRate),
		DataBits:              8,
		StopBits:              1,
		InterCharacterTimeout: 200,
		MinimumReadSize:       0,
	}
	port, err := serial.Open(options)
	if err != nil {
		panic(err)
	}
	defer port.Close()

	this.log.Debug("serial connect success!")

	result := []byte{}
	for {
		buf := make([]byte, 128)
		n, err := port.Read(buf)

		if err != nil {
			if err != io.EOF {
				panic(err)
			} else {
				break
			}
		} else if n != 0 {
			buf = buf[0:n]
			this.log.Debug("get serial data %v %v", n, buf)
			endIndex := -1
			for i := 0; i != n; i++ {
				if buf[i] == '\n' {
					endIndex = i
					break
				}
			}
			if endIndex != -1 {
				result = append(result, buf[0:endIndex]...)
				var data string
				if len(result) != 0 && result[len(result)-1] == 13 {
					data = string(result[0 : len(result)-1])
				} else {
					data = string(result)
				}
				if len(data) != 0 {
					this.msgChan <- data
				}
				if endIndex+1 >= len(buf) {
					result = []byte{}
				} else {
					result = buf[endIndex+1 : len(buf)]
				}
			} else {
				result = append(result, buf...)
			}
		}
	}
}

func (this *Serial) run() {
	for {
		func() {
			defer CatchCrash(func(e Exception) {
				this.log.Critical("Serial crash , will retry in next 20s , error:%v", e.Error())
				time.Sleep(time.Second * 20)
			})
			this.singleRun()
			this.log.Debug("serial disconnect!,will reconnect in 10s")
			time.Sleep(time.Second * 10)
		}()
	}
}

func (this *Serial) run_test() {
	time.Sleep(time.Second * 1)
	this.msgChan <- "ATF 0208091115018749401"
	time.Sleep(time.Second * 1)
	this.msgChan <- "ATF 0208091115018749402"
	time.Sleep(time.Second * 1)
	this.msgChan <- "ATRING"
	time.Sleep(time.Second * 1)
	this.msgChan <- "ATRING"
	time.Sleep(time.Second * 1)
	this.msgChan <- "ATF 0208091115018749403"
	time.Sleep(time.Second * 1)
	this.msgChan <- "ATRING"
	time.Sleep(time.Second * 1)
	this.msgChan <- "ATRING"
}

func (this *Serial) Run() error {
	go this.run()
	this.fireListener()
	return nil
}

func (this *Serial) Close() {

}
