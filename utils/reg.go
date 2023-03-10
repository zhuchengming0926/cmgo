/**
 * @File: reg.go
 * @Author: zhuchengming
 * @Description:正则表达式使用
 * @Date: 2021/5/10 15:47
 */

package utils

import (
	"fmt"
	"regexp"
)

func IsIP(str string) bool {
	regular := `((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsEmail(str string) bool {
	regular := `^[_a-zA-Z0-9-]+(\.[_a-zA-Z0-9-]+)*@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*(\.[a-zA-Z]{2,4})$`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsTelephone(str string) bool {
	regular := `^(010|02\d{1}|0[3-9]\d{2})-\d{7,9}(-\d+)?$`
	return regexp.MustCompile(regular).MatchString(str)
}

func Is400(str string) bool {
	regular := `^400(-\d{3,4}){2}$`
	return regexp.MustCompile(regular).MatchString(str)

}

func IsPhone(str string) bool {
	regular := `^(\+?86-?)?(18|15|13)[0-9]{9}$`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsYMD(str string) bool {
	//(?!0000)  闰年:2016-02-29
	regular := `^([0-9]{4}-((0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-8])|(0[13-9]|1[0-2])-(29|30)|(0[13578]|1[02])-31)|([0-9]{2}(0[48]|[2468][048]|[13579][26])|(0[48]|[2468][048]|[13579][26])00)-02-29)$`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsHMS_APM(str string) bool {
	//hh:mm:ss xx
	regular := `(0[1-9]|1[0-2]):[0-5][0-9]:[0-5][0-9] ([AP]M)`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsHMS(str string) bool {
	//hh:mm:ss
	regular := `(0[1-9]|1[0-2]):[0-5][0-9]:[0-5][0-9]`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsYMDHMS(str string) bool {
	//YYYY-MM-DD
	//(?!0000)  闰年:2016-02-29
	regular := `^([0-9]{4}-((0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-8])|(0[13-9]|1[0-2])-(29|30)|(0[13578]|1[02])-31)|([0-9]{2}(0[48]|[2468][048]|[13579][26])|(0[48]|[2468][048]|[13579][26])00)-02-29) (0[1-9]|1[0-2]):[0-5][0-9]:[0-5][0-9]$`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsNumber(str string) bool {
	//只能输入数字
	regular := `^[0-9]*$`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsFloat(str string) bool {
	//整数或者小数
	regular := `^[0-9]+([.]{0,1}[0-9]+){0,1}$`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsNumber_M_N(str string, m, n int) bool {
	//只能输入m-n个数字
	regular := fmt.Sprintf("^\\d{%d,%d)$", m, n)
	return regexp.MustCompile(regular).MatchString(str)
}

func IsSpecialSymbols(str string) bool {
	//特殊符号开头
	regular := `^[!@#$%^&*()_+-={}|[]\:";'<>?,./]+`
	return regexp.MustCompile(regular).MatchString(str)
}

func IsChineseCharacter(str string) bool {
	//是否全是中文
	regular := `^[\u4E00-\u9FA5]+$`
	return regexp.MustCompile(regular).MatchString(str)
}

//(?i)表示所在位置右侧的表达式开启忽略大小写模式
//https://blog.csdn.net/panamera918/article/details/80077170
func IsValidateImageUrl(url string) bool {
	imageCompile := regexp.MustCompile(`(?i)^(http|https)://.*?/[\da-z_\-]+\.(jpg|png)$`)
	name := imageCompile.FindString(url)
	if len(name) == 0 {
		return false
	}
	return true
}

func IsValidateVideoUrl(url string) bool {
	videoCompile := regexp.MustCompile(`(?i)^(http|https)://.*?/[\da-z_\-]+\.mp4$`)
	name := videoCompile.FindString(url)
	if len(name) == 0 {
		return false
	}
	return true
}

func IsValidateAudioUrl(url string) bool {
	videoCompile := regexp.MustCompile(`(?i)^(http|https)://.*?/[\da-z_\-]+\.mp3$`)
	name := videoCompile.FindString(url)
	if len(name) == 0 {
		return false
	}
	return true
}

func IsValidateName(name string) bool {
	nameCompile := regexp.MustCompile("(?i)^[\u4e00-\u9fa5a-z0-9,?!'，？！’ ‘]*?$")
	nameMatch := nameCompile.FindString(name)
	if len(nameMatch) == 0 {
		return false
	}
	return true
}