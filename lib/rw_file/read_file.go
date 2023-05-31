package rw_file

import (
	"bufio"
	"cmgo/utils"
	"io"
	"log"
	"os"
)

/**
 * @Author: chengming1
 * @Date: 2023/2/2 15:47
 * @Desc: 参考：https://blog.csdn.net/raoxiaoya/article/details/117998066
 */

// 1、直接读取整个文件(最优）
func ReadWholeFile(fileName string) string {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("os.ReadFile %s failed, err: %s", fileName, err.Error())
	}
	return utils.BytesToString(fileContent)
}

// 2、先从文件读取到file中，在从file读取到buf, buf在追加到最终的[]byte
func ReadToBufByte(fileName string) string {
	// 获取文件指针
	filePointer, err := os.Open(fileName)
	if err != nil {
		log.Printf("os.Open %s failed, err: %s", fileName, err.Error())
		return ""
	}

	// 读取文件内容到缓冲[]byte中
	defer filePointer.Close()
	var result []byte
	partContent := make([]byte, 1024)

	for {
		bytesLen, err := filePointer.Read(partContent)
		if err != nil && err != io.EOF {
			log.Printf("filePointer.Read %s failed, err: %s", fileName, err.Error())
			return ""
		}
		if bytesLen == 0 { //说明读取结束
			break
		}
		result = append(result, partContent[:bytesLen]...)
	}
	return utils.BytesToString(result)
}

// 3、先从文件读取到file, 在从file读取到Reader中，从Reader读取到buf, buf最终追加到[]byte
func ReadToReader(fileName string) string {
	filePointer, err := os.Open(fileName)
	if err != nil {
		log.Printf("os.Open %s failed, err: %s", fileName, err.Error())
		return ""
	}
	defer filePointer.Close()

	r := bufio.NewReader(filePointer)
	var result []byte
	partContent := make([]byte, 1024)
	for {
		bytesLen, err := r.Read(partContent)
		if err != nil && err != io.EOF {
			log.Printf("Reader.Read %s failed, err: %s", fileName, err.Error())
			return ""
		}
		if bytesLen == 0 { //说明读取结束
			break
		}
		result = append(result, partContent[:bytesLen]...)
	}
	return utils.BytesToString(result)
}

// 4、按行读取
func ReadLine(fileName string, callBack func(string) string) string {
	filePointer, err := os.Open(fileName)
	if err != nil {
		log.Printf("os.Open %s failed, err: %s", fileName, err.Error())
		return ""
	}
	defer filePointer.Close()

	// 以这个文件为参数，创建一个 scanner
	s := bufio.NewScanner(filePointer)

	// 扫描每行文件，按行读取
	for s.Scan() {
		log.Println(callBack(s.Text()))
	}
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
	return ""
}
