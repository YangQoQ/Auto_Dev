package Struct

import (
	"WeChatAutoDev/Dev/Util"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// 配置信息
type BaseConfig struct {
	// 数据库配置结构体
	Database struct {
		Ip             string `json:"ip"`             // 体检数据服务器
		Name           string `json:"name"`           // 账号
		Pwd            string `json:"pwd"`            // 密码
		Nqpeisname     string `json:"nqpeisname"`     // 体检数据库名称
		Wechatpeisname string `json:"wechatpeisname"` // 体检微信数据库名称
	} `json:"database"`
	// 文件部署地址 （拷贝到前置机上准备部署的文件路径）
	Devfilepath struct {
		Wechatapi          string `json:"wechatapi"`          // 微信体检服务Api部署路径
		Wechatworkplanapi  string `json:"wechatworkplanapi"`  // 微信排班系统Api部署路径
		Taskapi            string `json:"taskapi"`            // TaskApi部署文件路径
		Wechatview         string `json:"wechatview"`         // 微信体检前端页面部署路径
		Wechatworkplanview string `json:"wechatworkplanview"` // 微信排班页面部署路径
	} `json:"devfilepath"`
	// 文件存放地址 （部署文件存放路径）
	Storefilepath struct {
		Storedisk         string `json:"storedisk"`         // 磁盘存放位置
		Storewechatapi    string `json:"storewechatapi"`    // 微信体检服务api磁盘存放位置
		Storeworkplanapi  string `json:"storeworkplanapi"`  // 微信排班api磁盘存放位置
		Storetaskapi      string `json:"storetaskapi"`      // 微信Taskapi磁盘存放位置
		Storewechatview   string `json:"storewechatview"`   // 微信体检服务前端页面磁盘存放位置
		Stroeworkplanview string `json:"stroeworkplanview"` // 微信体检服务排班页面磁盘存放位置
	} `json:"storefilepath"`
}

func ReadFile() BaseConfig {
	// 读取配置文件
	config, err := ReadConfig("Dev/config.json")
	if err != nil {
		Util.ShowMessage("Error", "无法读取配置文件")
	}
	return config
}

// 读取config.json配置文件并适配Config对象
func ReadConfig(filename string) (BaseConfig, error) {
	var config BaseConfig

	// 打开并读取文件
	file, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("打开文件时发生错误: %v", err)
	}
	defer file.Close()

	// 检查文件大小
	stat, err := file.Stat()
	if err != nil {
		return config, fmt.Errorf("获取文件信息时发生错误: %v", err)
	}
	if stat.Size() == 0 {
		return config, errors.New("文件为空")
	}

	// 读取文件内容
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("读取文件时发生错误: %v", err)
	}

	// 解析 JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("解析 JSON 时发生错误: %v", err)
	}

	return config, nil
}
