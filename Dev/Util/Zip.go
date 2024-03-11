package Util

import (
	"archive/zip"
	"github.com/axgle/mahonia"
	"io"
	"os"
	"path/filepath"
)

// 程序包解压缩
func Unzip(src string, destDir string) (string, error) {
	// 第一步，打开 zip 文件
	zipFile, err := zip.OpenReader(src)
	if err != nil {
		ShowMessage("Error", "压缩文件格式必须为zip")
		panic(err)
	}
	defer zipFile.Close()

	filePath := mahonia.NewDecoder("gbk").ConvertString(filepath.Join(destDir, zipFile.File[0].Name))

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
	return filePath, nil
}
