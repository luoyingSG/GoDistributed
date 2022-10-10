package log

import (
	"fmt"
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type fileLog string

// 使得 fileLog 类型实现 io.Writer 接口定义的 Write 方法
func (fl fileLog) Write(data []byte) (int, error) {
	// 打开日志文件
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if nil != err {
		fmt.Printf("Failed to open log file %s\n", string(fl))
		return 0, err
	}
	// 关闭前关闭日志文件
	defer f.Close()
	// 写入日志
	return f.Write(data)
}

// 初始化日志系统
func Run(destination string) {
	// 初始化传入：io.Writer 类型的值，日志的前缀，日志的标记（？）
	log = stlog.New(fileLog(destination), "go ", stlog.LstdFlags)
}

// 一个请求处理函数
func RegisterHandlers() {
	// http 请求处理
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		// 处理 POST 请求
		case http.MethodPost:
			msg, err := ioutil.ReadAll(r.Body)
			if nil != err || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

// 将网络日志写入文件
func write(message string) {
	log.Printf("%v\n", message)
}
