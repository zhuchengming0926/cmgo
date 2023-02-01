package utils

import (
	"fmt"
	"testing"
)

/**
 * @Author: chengming1
 * @Date: 2023/1/30 下午2:35
 * @Desc:
 */

func TestGetLocalIp(t *testing.T) {
	localIp := GetLocalIp()
	fmt.Println(localIp)
}