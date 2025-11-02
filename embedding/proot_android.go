//go:build android

package embedding

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetProotPath 获取 proot 二进制文件路径
// proot 允许在 Android 上运行 Linux 二进制文件
func GetProotPath() (string, error) {
	libDir, err := findNativeLibraryDir()
	if err != nil {
		return "", fmt.Errorf("无法找到 lib 目录: %v", err)
	}

	// proot 文件名
	var prootName string
	switch runtime.GOARCH {
	case "arm64":
		prootName = "libproot.so"
	case "arm":
		prootName = "libproot.so"
	case "amd64":
		prootName = "libproot_x86_64.so"
	case "386":
		prootName = "libproot_x86.so"
	default:
		prootName = "libproot.so"
	}

	// 可能的子目录
	abiDirs := []string{"", "arm64-v8a", "arm64", "x86_64", "x86"}

	for _, abiDir := range abiDirs {
		baseDir := libDir
		if abiDir != "" {
			baseDir = filepath.Join(libDir, abiDir)
		}

		prootPath := filepath.Join(baseDir, prootName)
		if info, err := os.Stat(prootPath); err == nil && !info.IsDir() {
			// 确保有执行权限
			os.Chmod(prootPath, 0755)
			return prootPath, nil
		}
	}

	return "", fmt.Errorf("找不到 proot 二进制文件")
}

// SetupProotEnvironment 设置 proot 运行环境
func SetupProotEnvironment() (map[string]string, error) {
	env := make(map[string]string)

	// 获取应用的私有数据目录
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		// 尝试使用应用数据目录
		homeDir = "/data/data/com.oneclickvirt.goecs"
		if _, err := os.Stat(homeDir); os.IsNotExist(err) {
			// 如果无法访问，使用临时目录
			homeDir = "/data/local/tmp"
		}
	}

	env["HOME"] = homeDir
	env["TMPDIR"] = filepath.Join(homeDir, "tmp")
	env["PATH"] = "/usr/local/bin:/usr/bin:/bin:/usr/local/sbin:/usr/sbin:/sbin"

	// 创建必要的目录
	os.MkdirAll(env["TMPDIR"], 0755)

	return env, nil
}

// BuildProotCommand 构建使用 proot 运行的命令
// prootPath: proot 二进制文件路径
// ecsPath: ECS 二进制文件路径
// args: ECS 命令参数
func BuildProotCommand(prootPath, ecsPath string, args []string) (string, []string) {
	// proot 参数
	prootArgs := []string{
		"-0",      // 模拟 root 用户
		"-r", "/", // 根目录（不改变）
		"-w", "/", // 工作目录
		"--link2symlink", // 将硬链接转换为符号链接
	}

	// 添加 ECS 命令和参数
	prootArgs = append(prootArgs, ecsPath)
	prootArgs = append(prootArgs, args...)

	return prootPath, prootArgs
}
