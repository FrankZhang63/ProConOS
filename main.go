package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// Config banner结构体
type Config struct {
	LadderLogicRuntime string `json:"Ladder Logic Runtime"`
	PLCType            string `json:"PLC Type"`
	ProjectName        string `json:"Project Name"`
	BootProject        string `json:"Boot Project"`
	ProjectSourceCode  string `json:"Project Source Code"`
}

// 连接处理
func handleClient(conn net.Conn, response []string) {
	defer conn.Close()
	// 读取客户端请求的相关信息
	ipPort := strings.Split(conn.RemoteAddr().String(), ":")
	//请求时客户端发送的
	s := []byte{0xcc, 0x01, 0x00, 0x0b, 0x40, 0x02, 0x00, 0x00, 0x47, 0xee}
	fmt.Printf("[Proconos] 请求信息 IP:%s\tPort:%s\tType:%s\tMsg:%s\n", ipPort[0], ipPort[1], conn.RemoteAddr().Network(), string(s))

	// 把配置逐个变成字节切片
	LadderLogicRuntime := []byte(response[0])
	PLCType := []byte(response[1])
	ProjectName := []byte(response[2])
	BootProject := []byte(response[3])
	ProjectSourceCode := []byte(response[4])
	// 拼接返回数据
	data := []byte{0xcc, 0x01, 0x86, 0x04, 0x00, 0x02, 0x92, 0x00, 0x56, 0x34, 0x2e, 0x31}
	data = append(data, LadderLogicRuntime...)
	data = append(data, 0x00, 0x00, 0x00)
	data = append(data, PLCType...)
	data = append(data, 0x00, 0x00, 0x00, 0x00, 0x00)
	data = append(data, ProjectName...)
	data = append(data, 0x00)
	data = append(data, BootProject...)
	data = append(data, 0x00)
	data = append(data, ProjectSourceCode...)
	data = append(data, 0x00, 0x00, 0x45, 0x78, 0x69, 0x73, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x31, 0xa2, 0xc4, 0x61, 0x31, 0xa2, 0xc4, 0x61, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x5d, 0xff, 0xff, 0xff, 0xff, 0xbf, 0x0e, 0x00, 0x00, 0x44, 0xec, 0x03, 0x00, 0xf4, 0x18, 0x0f, 0x00, 0x00, 0x10, 0xe6, 0x03, 0xfa, 0x00, 0x01, 0x00, 0x00, 0x00, 0x47, 0x58)

	conn.Write(data)
}
func main() {
	//获取配置文件路径
	var configFilePath *string
	configFilePath = flag.String("config", "", "配置文件路径")
	flag.Parse()
	//读配置文件内容
	configFileContent, err := os.ReadFile(*configFilePath)
	if err != nil {
		fmt.Println("[Proconos] 配置文件指定协议相关配置错误")
		return
	}
	var config Config
	//反序列化到结构体
	if err := json.Unmarshal(configFileContent, &config); err != nil {
		return
	}
	//确定返回结构
	var response []string
	response = append(response, config.LadderLogicRuntime, config.PLCType, config.ProjectName, config.BootProject, config.ProjectSourceCode)
	// 创建监听Socket
	listener, err := net.Listen("tcp", ":20547")
	if err != nil {
		fmt.Println("[Proconos] 创建监听Socket时发生错误:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("[Proconos] 等待连接")
	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[Proconos] 接受连接时发生错误:", err)
			continue
		}
		// 启动一个协程来处理客户端连接
		go handleClient(conn, response)
	}
}
