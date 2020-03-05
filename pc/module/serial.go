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

	data := []byte{}
	dataLength := 0
	beginTime := time.Time{}
	beginCount := 0
	state := 0
	resetState := func() {
		data = nil
		dataLength = 0
		beginTime = time.Time{}
		beginCount = 0
		state = 0
	}

	handleMessage := func(data []byte) (bool, string) {
		timeStr := ""
		dataStr := ""
		if len(data) >= 10 &&
			data[0] == 0x01 &&
			data[1] == 0x08 {
			//时间数据
			timeStr = string(data[2:10])
			data = data[10:]
		} else {
			timeStr = time.Now().Format("01021504")
		}
		if len(data) >= 3 &&
			data[0] == 0x02 &&
			len(data) == int(data[1])+2 {
			dataStr = string(data[2:])
		} else {
			return false, ""
		}
		return true, "ATF " + timeStr + dataStr
	}

	handleChar := func(c byte) (bool, string) {
		if state == 0 {
			if c != 0x55 {
				//重置状态
				resetState()
			} else {
				//叠加状态
				beginCount++
				if beginCount > 10 {
					//足够数量的0x55时，转入状态1
					state = 1
					beginTime = time.Now()
				}
			}
		} else if state == 1 {
			if c == 0x80 {
				//一直等待一个字符为0x80
				state = 2
			}
		} else if state == 2 {
			//0x80之后的数字就是为长度
			dataLength = int(c)
			if dataLength <= 0 {
				resetState()
			} else {
				state = 3
			}
		} else if state == 3 {
			//获取具体数据
			data = append(data, c)
			if len(data) == dataLength {
				//足够数据了
				isValid, msg := handleMessage(data)
				resetState()
				if isValid {
					return true, msg
				}
			}
		}
		//默认都是没有数据的
		return false, ""
	}

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
			now := time.Now()
			if state != 0 &&
				now.Sub(beginTime) > time.Second*5 {
				//5秒内还没结束的，重置读取数据
				resetState()
			}

			for i := 0; i != n; i++ {
				//逐个处理字符
				hasData, msg := handleChar(buf[i])
				if hasData {
					this.log.Debug("get serial data %v %v", len(msg), msg)
					this.msgChan <- msg
				}
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
