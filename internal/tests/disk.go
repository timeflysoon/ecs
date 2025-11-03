package tests

import (
	"runtime"
	"strings"

	"github.com/oneclickvirt/disktest/disk"
)

func DiskTest(language, testMethod, testPath string, isMultiCheck bool, autoChange bool) (realTestMethod, res string) {
	// Android 平台不支持二进制文件测试
	if runtime.GOOS == "android" {
		realTestMethod = "disabled"
		if language == "en" {
			res = "Disk test is not supported on Android platform.\n" +
				"Reason: Android security sandbox prevents apps from executing\n" +
				"Alternative: Please use Termux to run the binary test from:\n" +
				"https://github.com/oneclickvirt/ecs\n"
		} else {
			res = "Android 平台不支持硬盘测试。\n" +
				"原因：Android 安全沙箱机制禁止应用直接执行外部二进制文件\n" +
				"替代方案：请使用 Termux 执行以下项目的二进制文件测试：\n" +
				"https://github.com/oneclickvirt/ecs\n"
		}
		return
	}

	switch testMethod {
	case "fio":
		res = disk.FioTest(language, isMultiCheck, testPath)
		if res == "" && autoChange {
			res += disk.DDTest(language, isMultiCheck, testPath)
			realTestMethod = "dd"
		} else {
			realTestMethod = "fio"
		}
	case "dd":
		res = disk.DDTest(language, isMultiCheck, testPath)
		if res == "" && autoChange {
			res += disk.FioTest(language, isMultiCheck, testPath)
			realTestMethod = "fio"
		} else {
			realTestMethod = "dd"
		}
	default:
		if runtime.GOOS == "windows" {
			realTestMethod = "winsat"
			res = disk.WinsatTest(language, isMultiCheck, testPath)
		} else {
			res = disk.DDTest(language, isMultiCheck, testPath)
			realTestMethod = "dd"
		}
	}
	if !strings.Contains(res, "\n") && res != "" {
		res += "\n"
	}
	return
}
