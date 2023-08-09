package rw_file

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

func GetDirAllFileContent() {
	folder := "/realtime_feature_compute" // 文件夹路径
	outputFile := "/all.txt"              // 输出文件路径

	// 创建一个io.Writer，用于写入文件内容
	writer, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writer.Close()
	w := bufio.NewWriter(writer)

	// 遍历文件夹下的所有文件
	filePaths, err := GetDirAllFilePaths(folder)
	if err != nil {
		fmt.Println(err)
	}

	for _, filePath := range filePaths {
		_, fileName := filepath.Split(filePath)
		// fmt.Println(fileName)
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("os.ReadFile %s failed, err: %s", filePath, err.Error())
		}
		cnt, err := w.WriteString("// " + fileName + "\n" + BytesToString(fileContent) + "\n" + "\n")
		fmt.Println(cnt, err)
	}

	fmt.Println(" finished!")
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 递归获取文件目录下所有文件path
func GetDirAllFilePaths(dirname string) ([]string, error) {
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))
	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetDirAllFilePaths(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}
