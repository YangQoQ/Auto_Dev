// 微信自动化部署脚本程序
package main

import (
	"WeChatAutoDev/Dev/Struct"
	"WeChatAutoDev/Dev/Util"
	"github.com/axgle/mahonia"
	"log"
	"os"
	"os/exec"
	"strings"
)

// 配置信息
type Config struct {
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

func main() {
	// 读取配置文件
	config := Struct.ReadFile()

	// 创建存放文件夹
	folderPaths := [...]string{
		config.Storefilepath.Storewechatapi,
		config.Storefilepath.Storeworkplanapi,
		config.Storefilepath.Storetaskapi,
		config.Storefilepath.Storewechatview,
		config.Storefilepath.Stroeworkplanview,
	}

	// 部署文件夹键值对
	devFolderMap := make(map[string]string)
	devFolderMap[config.Storefilepath.Storewechatapi] = config.Devfilepath.Wechatapi
	devFolderMap[config.Storefilepath.Storeworkplanapi] = config.Devfilepath.Wechatworkplanapi
	devFolderMap[config.Storefilepath.Storetaskapi] = config.Devfilepath.Taskapi
	devFolderMap[config.Storefilepath.Storewechatview] = config.Devfilepath.Wechatview
	devFolderMap[config.Storefilepath.Stroeworkplanview] = config.Devfilepath.Wechatworkplanview

	// 检查是否安装了iis
	iisCmdCheck := exec.Command("sc", "query", "w3svc", "STATE")
	iisout, _ := iisCmdCheck.CombinedOutput()
	iiResult := mahonia.NewDecoder("gbk").ConvertString(string(iisout))
	if !strings.Contains(iiResult, "RUNNING") {
		Util.ShowMessage("Error", "请安装iis")
	}

	// 检查是否安装了microsoftrootcertificateauthority2011微软证书
	cmd := exec.Command("certutil", "-store", "-user", "Root")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	// 输出字符编码转gbk
	cmdResult := mahonia.NewDecoder("gbk").ConvertString(string(out))
	// 是否存在根证书
	if strings.Contains(cmdResult, "受信任的根证书颁发机构") {
		// 检查服务器是否安装了 .NET Framework 4.7.2 或者以上版本
		ckFrameworkCmd := exec.Command("reg", "query", "HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v4\\Full", "/v", "Release")
		out, _ := ckFrameworkCmd.CombinedOutput()
		if strings.Contains(string(out), "0x461808") || strings.Contains(string(out), "0x82348") {
			// 检查是否存在配置存放盘符
			_, err := os.Stat(config.Storefilepath.Storedisk)
			// 存在盘符
			if err == nil {
				for _, filePath := range folderPaths {
					// 检查是否已经存在文件夹
					_, storeFilePath := os.Stat(filePath)
					if storeFilePath != nil {
						err := os.MkdirAll(filePath, os.ModePerm)
						if err != nil {
							Util.ShowMessage("Error", "文件夹创建失败"+filePath)
							break
						}

						// 没有找到对应文件夹Return
						if devFolderMap[filePath] == "" {
							return
						}

						// 解压文件到指定目录
						Util.Unzip(devFolderMap[filePath], filePath)

						// 修改各个项目配置文件中数据库连接，通过存放路径获取到webconfig文件

						// 部署iis

					}
				}

			} else if os.IsNotExist(err) {
				Util.ShowMessage("Error", "不存在磁盘"+config.Storefilepath.Storedisk)
			}
			return
		} else {
			// 安装 .NET Framework 4.7.2
			Util.ShowMessage("Error", "请安装 .NET Framework 4.7.2")
			// cmd := exec.Command("powershell", "Start-Process", "-Verb", "runAs", "D:\\ndp472-devpack-enu.exe")
			// err := cmd.Run()
			// if err != nil {// log.Fatal(err)
			// }
			return
		}
	} else {
		Util.ShowMessage("Error", "失败原因: 需安装微软根颁发机构的证书")
	}
}
