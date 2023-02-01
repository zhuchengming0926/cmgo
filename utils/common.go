package utils

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"os"
	"strings"
	"unsafe"
)

// 判断文件或目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

func LoadFile(path string) ([]byte, bool) {
	fileSize, err := FileSize(path)
	if nil != err || fileSize < 1 {
		return nil, false
	}
	buf := make([]byte, fileSize)
	f, err := os.Open(path)
	defer f.Close()

	if nil != err {
		return nil, false
	}
	n, err := f.Read(buf)
	if nil != err || fileSize != int64(n) {
		return nil, false
	}

	return buf, true
}

func LoadConf(path string, cf interface{}) bool {
	fileSize, err := FileSize(path)
	if nil != err || fileSize < 1 {
		return false
	}
	buf := make([]byte, fileSize)
	f, err := os.Open(path)
	defer f.Close()
	if nil != err {
		return false
	}
	n, err := f.Read(buf)
	if nil != err || fileSize != int64(n) {
		return false
	}
	if nil != UnmarshalJson(buf, cf) {
		return false
	}
	return true
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func SaveConf(conf interface{}, path string) bool {
	buf, err := jsoniter.MarshalIndent(conf, "", "    ")
	if nil != err {
		return false
	}
	f, err := os.Create(path)
	if nil != err {
		return false
	}
	defer f.Close()
	size := len(buf)
	n, err := f.Write(buf)
	if nil != err || n != size {
		return false
	}
	return true
}

func UnmarshalJson(jsonVal []byte, objVal interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(BytesToString(jsonVal)))
	decoder.UseNumber()
	err := decoder.Decode(objVal)
	return err
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}