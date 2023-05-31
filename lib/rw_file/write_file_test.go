package rw_file

import "testing"

/**
 * @Author: chengming1
 * @Date: 2023/2/3 09:35
 * @Desc:
 */

func TestBufioAppendWriteFile(t *testing.T) {
	err := BufioAppendWriteFile("./test.txt", "你好")
	if err != nil {
		t.Error(err)
	}
}

func TestCoverWriteFile(t *testing.T) {
	err := CoverWriteFile("./test1.txt", "你好")
	if err != nil {
		t.Error(err)
	}
}
