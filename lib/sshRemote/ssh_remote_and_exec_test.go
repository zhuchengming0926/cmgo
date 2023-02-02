package sshRemote

import "testing"
import "fmt"

func TestRunShell(t *testing.T) {
	cliConf := new(ClientConfig)
	cliConf.createClient("10.41.11.x", 22, "xxx", "xxx")
	/*
		可以看到我们这里每次执行一条命令都会创建一条session
		这是因为一条session默认只能执行一条命令
		并且两条命令不可以分开写
		比如：
		cliConf.RunShell("cd /opt")
		cliConf.RunShell("ls")
		这两条命令是无法连续的，下面的ls查看的依旧是~目录
		因此我们可以连着写，使用;分割
	*/
	fmt.Println(cliConf.RunShell("cd /opt; ls -l"))
	/*
		total 20
		drwxr-xr-x 3 root root 4096 Nov 18 14:05 hadoop
		drwxr-xr-x 3 root root 4096 Nov 18 14:20 hive
		drwxr-xr-x 3 root root 4096 Nov 18 15:07 java
		drwxr-xr-x 3 root root 4096 Nov  4 23:01 kafka
		drwxr-xr-x 3 root root 4096 Nov  4 22:54 zookeeper
	*/
}
