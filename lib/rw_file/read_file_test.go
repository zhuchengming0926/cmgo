package rw_file

import "testing"

/**
 * @Author: chengming1
 * @Date: 2023/2/3 09:54
 * @Desc:
 */

func TestReadLine(t *testing.T) {
	ReadLine("./test.txt", func(s string) string {
		return s + "坚强"
	})
}
