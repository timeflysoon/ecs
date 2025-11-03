package tests

import (
	"runtime"
	"strings"

	"github.com/oneclickvirt/memorytest/memory"
)

func MemoryTest(language, testMethod string) (realTestMethod, res string) {
	testMethod = strings.ToLower(testMethod)
	if testMethod == "" {
		testMethod = "auto"
	}

	// Android 平台不支持二进制文件测试
	if runtime.GOOS == "android" {
		realTestMethod = "disabled"
		if language == "en" {
			res = "Memory test is not supported on Android platform.\n" +
				"Reason: Android security sandbox prevents apps from executing\n" +
				"Alternative: Please use Termux to run the binary test from:\n" +
				"https://github.com/oneclickvirt/ecs\n"
		} else {
			res = "Android 平台不支持内存测试。\n" +
				"原因：Android 安全沙箱机制禁止应用直接执行外部二进制文件\n" +
				"替代方案：请使用 Termux 执行以下项目的二进制文件测试：\n" +
				"https://github.com/oneclickvirt/ecs\n"
		}
		return
	}

	if runtime.GOOS == "windows" {
		switch testMethod {
		case "stream":
			res = memory.WinsatTest(language)
			realTestMethod = "winsat"
		case "dd":
			res = memory.WindowsDDTest(language)
			if res == "" || strings.TrimSpace(res) == "" {
				res += memory.WinsatTest(language)
				realTestMethod = "winsat"
			} else {
				realTestMethod = "dd"
			}
		case "sysbench":
			res = memory.WinsatTest(language)
			realTestMethod = "winsat"
		case "auto", "winsat":
			res = memory.WinsatTest(language)
			realTestMethod = "winsat"
		default:
			res = memory.WinsatTest(language)
			realTestMethod = "winsat"
		}
	} else {
		switch testMethod {
		case "stream":
			res = memory.StreamTest(language)
			if res == "" || strings.TrimSpace(res) == "" {
				res += memory.DDTest(language)
				realTestMethod = "dd"
			} else {
				realTestMethod = "stream"
			}
		case "dd":
			res = memory.DDTest(language)
			realTestMethod = "dd"
		case "sysbench":
			res = memory.SysBenchTest(language)
			if res == "" || strings.TrimSpace(res) == "" {
				res += memory.DDTest(language)
				realTestMethod = "dd"
			} else {
				realTestMethod = "sysbench"
			}
		case "auto":
			res = memory.StreamTest(language)
			if res == "" || strings.TrimSpace(res) == "" {
				res = memory.DDTest(language)
				if res == "" || strings.TrimSpace(res) == "" {
					res = memory.SysBenchTest(language)
					if res == "" || strings.TrimSpace(res) == "" {
						realTestMethod = ""
					} else {
						realTestMethod = "sysbench"
					}
				} else {
					realTestMethod = "dd"
				}
			} else {
				realTestMethod = "stream"
			}
		case "winsat":
			// winsat 仅 Windows 支持，非 Windows fallback 到 dd
			res = memory.DDTest(language)
			realTestMethod = "dd"
		default:
			res = "Unsupported test method"
			realTestMethod = ""
		}
	}
	if !strings.Contains(res, "\n") && res != "" {
		res += "\n"
	}
	return
}
