package rw_file

import (
	"bufio"
	"log"
	"os"
)

/**
 * @Author: chengming1
 * @Date: 2023/2/2 15:48
 * @Desc:
 */

func CoverWriteFile(fileName string, content string) error {
	filePointer, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer filePointer.Close()

	w := bufio.NewWriter(filePointer) //创建新的 Writer 对象
	cnt, err := w.WriteString(content + "\n")
	if err != nil {
		return err
	}
	log.Printf("覆盖写入%d字节\n", cnt)
	w.Flush()
	return nil
}

func BufioAppendWriteFile(fileName string, content string) error {
	filePointer, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer filePointer.Close()

	w := bufio.NewWriter(filePointer) //创建新的 Writer 对象
	cnt, err := w.WriteString(content + "\n")
	if err != nil {
		return err
	}
	log.Printf("追加写入%d字节\n", cnt)
	w.Flush()
	return nil
}

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
