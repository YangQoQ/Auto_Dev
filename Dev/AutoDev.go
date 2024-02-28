// 微信自动化部署脚本程序
package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
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
	config, err := ReadConfig("Dev/config.json")
	if err != nil {
		ShowMessage("Error", "无法读取配置文件")
		return
	}

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
		ShowMessage("Error", "请安装iis")
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
							ShowMessage("Error", "文件夹创建失败"+filePath)
							break
						}

						if devFolderMap[filePath] == "" {
							return
						}

						// 解压文件到指定目录
						Unzip(devFolderMap[filePath], filePath)

						// 修改各个项目配置文件中数据库连接，通过存放路径获取到webconfig文件

						// 部署iis

					}
				}

			} else if os.IsNotExist(err) {
				ShowMessage("Error", "不存在磁盘"+config.Storefilepath.Storedisk)
			}
			return
		} else {
			// 安装 .NET Framework 4.7.2
			ShowMessage("Error", "请安装 .NET Framework 4.7.2")
			// cmd := exec.Command("powershell", "Start-Process", "-Verb", "runAs", "D:\\ndp472-devpack-enu.exe")
			// err := cmd.Run()
			// if err != nil {// log.Fatal(err)
			// }
			return
		}
	} else {
		ShowMessage("Error", "失败原因: 需安装微信任根颁发机构的证书")
	}
}

// 弹框提示组成方法
func IntPtr(n int) uintptr {
	return uintptr(n)
}

// 弹框提示组成方法
func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

// ShowMessage windows下的另一种DLL方法调用
func ShowMessage(tittle, text string) {
	user32dll, _ := syscall.LoadLibrary("user32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBoxW := user32.NewProc("MessageBoxW")
	MessageBoxW.Call(IntPtr(0), StrPtr(text), StrPtr(tittle), IntPtr(0))
	defer syscall.FreeLibrary(user32dll)
}

// 读取config.json配置文件并适配Config对象
func ReadConfig(filename string) (Config, error) {
	var config Config

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

// 程序包解压缩
func Unzip(src string, destDir string) error {
	// 第一步，打开 zip 文件
	zipFile, err := zip.OpenReader(src)
	if err != nil {
		ShowMessage("Error", "压缩文件格式必须为zip")
		panic(err)
	}
	defer zipFile.Close()

	// 第二步，遍历 zip 中的文件
	for _, f := range zipFile.File {
		filePath := mahonia.NewDecoder("gbk").ConvertString(filepath.Join(destDir, f.Name))
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		// 创建对应文件夹
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}
		// 解压到的目标文件
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}
		file, err := f.Open()
		if err != nil {
			panic(err)
		}
		// 写入到解压到的目标文件
		if _, err := io.Copy(dstFile, file); err != nil {
			panic(err)
		}
		dstFile.Close()
		file.Close()
	}
	return nil
}
