package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var TimeLocation *time.Location
var TimeFormat = "2006-01-02 15:04:05"
var DateFormat = "2006-01-02"
var LocalIP = net.ParseIP("127.0.0.1")

//公共初始化函数：支持两种方式设置配置文件
//
//函数传入配置文件 Init("./conf/dev/")
//如果配置文件为空，会从命令行中读取 	  -config conf/dev/
func Init(configPath string) error {
	return InitModule(configPath, []string{"base", "mysql", "redis"})
}

//模块初始化
func InitModule(configPath string, modules []string) error {
	conf := flag.String("config", configPath, "input config file like ./conf/dev/")
	flag.Parse()
	if *conf == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO]  config=%s\n", *conf)
	log.Printf("[INFO] %s\n", " start loading resources.")

	// 设置ip信息，优先设置便于日志打印
	ips := GetLocalIPs()
	if len(ips) > 0 {
		LocalIP = ips[0]
	}

	// 解析配置文件目录
	if err := ParseConfPath(*conf); err != nil {
		return err
	}

	//初始化配置文件
	if err := InitViperConf(); err != nil {
		return err
	}

	// 加载base配置
	if InArrayString("base", modules) {
		if err := InitBaseConf(GetConfPath("base")); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitBaseConf:"+err.Error())
		}
	}

	// 设置时区
	if location, err := time.LoadLocation(ConfBase.TimeLocation); err != nil {
		return err
	} else {
		TimeLocation = location
	}

	log.Printf("[INFO] %s\n", " success loading resources.")
	log.Println("------------------------------------------------------------------------")
	return nil
}

func NewTrace() *TraceContext {
	trace := &TraceContext{}
	trace.TraceId = GetTraceId()
	trace.SpanId = NewSpanId()
	return trace
}

func NewSpanId() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}

func GetTraceId() (traceId string) {
	return calcTraceId(LocalIP.String())
}

func calcTraceId(ip string) (traceId string) {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIP := net.ParseIP(ip)
	if netIP == nil {
		b.WriteString("00000000")
	} else {
		b.WriteString(hex.EncodeToString(netIP.To4()))
	}
	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0") // 末两位标记来源,b0为go

	return b.String()
}

func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return ips
}

func InArrayString(s string, arr []string) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

//Substr 字符串的截取
func Substr(str string, start int64, end int64) string {
	length := int64(len(str))
	if start < 0 || start > length {
		return ""
	}
	if end < 0 {
		return ""
	}
	if end > length {
		end = length
	}
	return string(str[start:end])
}

func ParseServerAddr(serverAddr string) (host, port string) {
	serverInfo := strings.Split(serverAddr, ":")
	if len(serverInfo) == 2 {
		host = serverInfo[0]
		port = serverInfo[1]
	} else {
		host = serverAddr
		port = ""
	}
	return host, port
}

func addTrace2Header(request *http.Request, trace *TraceContext) *http.Request {
	traceId := trace.TraceId
	cSpanId := NewSpanId()
	if traceId != "" {
		request.Header.Set("didi-header-rid", traceId)
	}
	if cSpanId != "" {
		request.Header.Set("didi-header-spanid", cSpanId)
	}
	trace.CSpanId = cSpanId
	return request
}

func GetMd5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encode(data string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func AddGetDataToUrl(urlString string, data url.Values) string {
	if strings.Contains(urlString, "?") {
		urlString = urlString + "&"
	} else {
		urlString = urlString + "?"
	}
	return fmt.Sprintf("%s%s", urlString, data.Encode())
}
