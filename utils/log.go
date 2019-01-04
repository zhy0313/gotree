package utils

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

const (
	LOG_ERROR = iota //错误日志
	LOG_WARN  = iota //警告日志
	LOG_INFO  = iota //普通日志
)

const (
	LOG_ERROR_DIR_NAME = "error" //错误日志目录名
	LOG_WARN_DIR_NAME  = "warn"  //警告日志目录名
	LOG_INFO_DIR_NAME  = "info"  //普通日志目录名
)

type log struct {
	errorDir string
	warnDir  string
	infoDir  string

	errFile  *os.File
	warnFile *os.File
	infoFile *os.File

	errFileName  string
	warnFileName string
	infoFileName string

	debug    bool
	mutexErr *sync.Mutex
	infoMsg  chan string
	warnMsg  chan string
	closeMsg chan bool

	stream  chan *logItem
	bseqMap bseqInterface
}

type logItem struct {
	category int
	data     []interface{}
}

type bseqInterface interface {
	Get(key string) interface{}
}

//Create 创建本结构
func (self *log) Init(dir string, bdict bseqInterface) {
	self.bseqMap = bdict
	logPath := dir
	if !PathExists(logPath) {
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	self.errorDir = logPath + "/" + LOG_ERROR_DIR_NAME
	if !PathExists(self.errorDir) {
		os.Mkdir(self.errorDir, os.ModePerm)
	}

	self.warnDir = logPath + "/" + LOG_WARN_DIR_NAME
	if !PathExists(self.warnDir) {
		os.Mkdir(self.warnDir, os.ModePerm)
	}

	self.infoDir = logPath + "/" + LOG_INFO_DIR_NAME
	if !PathExists(self.infoDir) {
		os.Mkdir(self.infoDir, os.ModePerm)
	}

	self.mutexErr = new(sync.Mutex)
	self.warnMsg = make(chan string, Config().DefaultInt("sys::LogWarnQueueLen", 512))
	self.infoMsg = make(chan string, Config().DefaultInt("sys::LogInfoQueueLen", 2048))
	self.closeMsg = make(chan bool, 2)
	go self.warnRun()
	go self.infoRun()
}

//Run
func (self *log) warnRun() {
	for {
		msg := <-self.warnMsg
		if msg == "" {
			self.closeMsg <- true
			break
		}
		self.write(LOG_WARN, msg)
	}
}

//Run
func (self *log) infoRun() {
	for {
		msg := <-self.infoMsg
		if msg == "" {
			self.closeMsg <- true
			break
		}
		self.write(LOG_INFO, msg)
	}
}

//WriteError 写入错误
func (self *log) WriteError(str ...interface{}) {
	stack := string(debug.Stack())
	str = append(str, "stack:"+stack)
	text := []interface{}{}
	if no := self.BseqNo(); no != "" {
		text = append(text, "bseq:"+no)
	}
	text = append(text, str...)
	datetime := self.nowDateTime()
	if self.debug {
		fmt.Println(Red(datetime + " " + fmt.Sprint(text)))
	}
	defer self.mutexErr.Unlock()
	self.mutexErr.Lock()
	self.write(LOG_ERROR, datetime+" "+fmt.Sprint(text))
}

func (self *log) WriteDaemonError(str ...interface{}) {
	text := []interface{}{}
	datetime := self.nowDateTime()
	text = append(text, str...)
	self.write(LOG_ERROR, datetime+" "+fmt.Sprint(text))
}

//WriteError 写入警告
func (self *log) WriteWarn(str ...interface{}) {
	text := []interface{}{}
	if no := self.BseqNo(); no != "" {
		text = append(text, "bseq:"+no)
	}
	datetime := self.nowDateTime()
	text = append(text, str...)
	if self.debug {
		fmt.Println(Red(datetime + " " + fmt.Sprint(text)))
	}

	self.warnMsg <- datetime + " " + fmt.Sprint(text)
}

//WriteInfo 写入普通
func (self *log) WriteInfo(str ...interface{}) {
	text := []interface{}{}
	if no := self.BseqNo(); no != "" {
		text = append(text, "bseq:"+no)
	}
	datetime := self.nowDateTime()
	text = append(text, str...)
	if self.debug {
		fmt.Println(Green(datetime + " " + fmt.Sprint(text)))
	}
	self.infoMsg <- datetime + " " + fmt.Sprint(text)
}

//WriteDebug 等同WriteInfo,只在dev模式下生效
func (self *log) WriteDebug(str ...interface{}) {
	if !self.debug {
		return
	}
	self.WriteInfo(str...)
}

//Write 写入日志
func (self *log) write(logType int, str string) {
	wt := self.getLogger(logType)
	if wt == nil {
		return
	}

	buf := []byte(str)
	wt.Write(append(buf, '\n'))
}

func (self *log) Debug() {
	self.debug = true
}

//Close 关闭日志
func (self *log) Close() {
	self.infoMsg <- ""
	self.warnMsg <- ""

	closeLen := 0
	for index := 0; index < 20; index++ {
		if closeLen == 2 {
			break
		}
		select {
		case _ = <-self.closeMsg:
			closeLen += 1
		default:
			//最多等待2秒
			time.Sleep(time.Millisecond * 100)
			continue
		}
	}

	self.errFile.Close()
	self.infoFile.Close()
	self.warnFile.Close()
}

//获取logger
func (self *log) getLogger(logType int) *os.File {
	var result *os.File

	fileName := self.getFileName(logType)
	switch logType {
	case LOG_ERROR:
		if self.errFile == nil {
			self.errFile, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0660)
			self.errFileName = fileName
		} else {
			if self.errFileName != fileName {
				self.errFile.Close()
				self.errFile, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0660)
				self.errFileName = fileName
			}
		}
		result = self.errFile
	case LOG_WARN:
		if self.warnFile == nil {
			self.warnFile, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0660)
			self.warnFileName = fileName
		} else {
			if self.warnFileName != fileName {
				self.warnFile.Close()
				self.warnFile, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0660)
				self.warnFileName = fileName
			}
		}
		result = self.warnFile
	default:
		if self.infoFile == nil {
			self.infoFile, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0660)
			self.infoFileName = fileName
		} else {
			if self.infoFileName != fileName {
				self.infoFile.Close()
				self.infoFile, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0660)
				self.infoFileName = fileName
			}
		}
		result = self.infoFile
	}

	return result
}

//getFileName 获取文件名
func (self *log) getFileName(logType int) string {
	dateDir := self.getDate(logType)
	hour, _, _ := time.Now().Clock()
	fileName := dateDir + "/" + strconv.Itoa(hour) + ".log"
	if PathExists(fileName) {
		return fileName
	}
	file, _ := os.Create(fileName)
	file.Close()
	return fileName
}

//GetDate 获取当前日期
func (self *log) getDate(logType int) string {
	var path string
	switch logType {
	case LOG_ERROR:
		path = self.errorDir
	case LOG_WARN:
		path = self.warnDir
	default:
		path = self.infoDir
	}
	year, month, day := time.Now().Date()

	dateDir := path + "/" + fmt.Sprintf("%d-%d-%d", year, month, day)
	if !PathExists(dateDir) {
		os.Mkdir(dateDir, os.ModePerm)
	}
	return dateDir
}

//NowDateTime 当前日期和时间
func (self *log) nowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (self *log) BseqNo() string {
	if self.bseqMap == nil {
		return ""
	}

	bseq := self.bseqMap.Get("bseq")
	if bseq == nil {
		return ""
	}

	str, ok := bseq.(string)
	if !ok {
		return ""
	}
	return str
}

var _log *log

//获取日志单例类
func Log() *log {
	if _log == nil {
		_log = new(log)
	}
	return _log
}

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func Red(str string) string {
	return textColor(TextRed, str)
}

func Green(str string) string {
	return textColor(TextGreen, str)
}

func Yellow(str string) string {
	return textColor(TextYellow, str)
}

func Blue(str string) string {
	return textColor(TextBlue, str)
}

func Magenta(str string) string {
	return textColor(TextMagenta, str)
}

func Cyan(str string) string {
	return textColor(TextCyan, str)
}

func White(str string) string {
	return textColor(TextWhite, str)
}

func Black(str string) string {
	return textColor(TextBlack, str)
}

func textColor(color int, str string) string {
	switch color {
	case TextBlack:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlack, str)
	case TextRed:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextRed, str)
	case TextGreen:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextGreen, str)
	case TextYellow:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextYellow, str)
	case TextBlue:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlue, str)
	case TextMagenta:
		return fmt.Sprintf("\\x1b[0;%dm%s\x1b[0m", TextMagenta, str)
	case TextCyan:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextCyan, str)
	case TextWhite:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextWhite, str)
	default:
		return str
	}
}

//PathExists 判断文件夹或文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
