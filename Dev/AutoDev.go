// 微信自动化部署脚本程序
package main

import (
	"WeChatAutoDev/Dev/Struct"
	"WeChatAutoDev/Dev/Util"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/axgle/mahonia"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
						destPath, _ := Util.Unzip(devFolderMap[filePath], filePath)

						// 修改各个项目配置文件中数据库连接，通过存放路径获取到webconfig
						targetFile, err := HandelWebConfig(destPath)
						if err != nil {
							Util.ShowMessage("Error", err.Error())
							return
						}

						// 找到文件后修改文件中DB配置
						fmt.Println("找到Web.Config文件:", targetFile)

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

// 处理WebConfig文件
func HandelWebConfig(filePath string) (string, error) {
	// 找到对应的配置文件
	targetFile := filepath.Join(filePath, "Web.config")
	// 检查文件是否存在
	_, err := os.Stat(targetFile)
	if err == nil {
		// 找到文件后修改配置文件
		_, err := ModifyWebConfig(targetFile)
		if err != nil {
			return "修改WebConfig文件失败", err
		}
	} else if os.IsNotExist(err) {
		return "", fmt.Errorf("文件 %s 不存在", targetFile)
	}
	return "", err
}

// 修改WebConfig指定结点
func ModifyWebConfig(modeifyPath string) (string, error) {
	// 价值Xml文件
	webXMl, err := os.Open(modeifyPath)
	if err != nil {
		panic(err)
	}

	// Parse XML document.
	weDoc, err := xmlquery.Parse(webXMl)
	if err != nil {
		panic(err)
	}

	// Find the connection string node
	node := xmlquery.FindOne(weDoc, "//connectionStrings//add[@name='SULL3']")
	if node != nil {
		// Modify the connection string attribute
		node.Attr[1].Value = "YourNewConnectionString"
		// 结点修改后保存文件
		
	} else {
		fmt.Println("Connection string node not found.")
	}
	return "", nil
}
