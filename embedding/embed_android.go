//go:build android

package embedding

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// findNativeLibraryDir 查找应用的 native library 目录
// 这个目录是系统自动管理的，包含从 APK 中提取的 .so 文件
func findNativeLibraryDir() (string, error) {
	// 方法 1: 通过环境变量获取（Fyne 可能会设置）
	if libDir := os.Getenv("ANDROID_LIB_DIR"); libDir != "" {
		if info, err := os.Stat(libDir); err == nil && info.IsDir() {
			return libDir, nil
		}
	}

	// 方法 2: 通过可执行文件路径推断
	execPath, err := os.Executable()
	if err == nil {
		// 可执行文件通常在 /data/app/<package>-<hash>/base.apk 或 /data/app/<package>-<hash>/oat/arm64/base.odex
		// native library 通常在 /data/app/<package>-<hash>/lib/arm64/

		// 尝试找到应用根目录
		dir := execPath
		for i := 0; i < 10; i++ { // 最多向上查找10层
			dir = filepath.Dir(dir)
			if dir == "/" || dir == "." {
				break
			}

			libDir := filepath.Join(dir, "lib")

			// 检查 lib 目录
			if info, err := os.Stat(libDir); err == nil && info.IsDir() {
				// 检查是否包含架构子目录或 .so 文件
				entries, err := os.ReadDir(libDir)
				if err == nil && len(entries) > 0 {
					return libDir, nil
				}
			}
		}
	}

	// 方法 3: 尝试标准的 Android native library 路径
	possibleBasePaths := []string{
		"/data/data/com.oneclickvirt.goecs/lib",
		"/data/app/com.oneclickvirt.goecs/lib",
	}

	for _, basePath := range possibleBasePaths {
		if info, err := os.Stat(basePath); err == nil && info.IsDir() {
			return basePath, nil
		}

		// 尝试带哈希的路径（Android 5.0+）
		parent := filepath.Dir(basePath)
		parentEntries, err := os.ReadDir(parent)
		if err == nil {
			for _, entry := range parentEntries {
				if entry.IsDir() && strings.HasPrefix(entry.Name(), "com.oneclickvirt.goecs") {
					libDir := filepath.Join(parent, entry.Name(), "lib")
					if info, err := os.Stat(libDir); err == nil && info.IsDir() {
						return libDir, nil
					}
				}
			}
		}
	}

	// 方法 4: 搜索 /data/app 目录
	dataAppDir := "/data/app"
	if entries, err := os.ReadDir(dataAppDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() && strings.Contains(entry.Name(), "com.oneclickvirt.goecs") {
				libDir := filepath.Join(dataAppDir, entry.Name(), "lib")
				if info, err := os.Stat(libDir); err == nil && info.IsDir() {
					return libDir, nil
				}
			}
		}
	}

	// 方法 5: 尝试通过 /proc/self/maps 查找已加载的共享库路径
	if mapsData, err := os.ReadFile("/proc/self/maps"); err == nil {
		lines := strings.Split(string(mapsData), "\n")
		for _, line := range lines {
			// 查找包含 .so 的行
			if strings.Contains(line, ".so") && strings.Contains(line, "/data/") {
				// 提取路径部分
				parts := strings.Fields(line)
				if len(parts) >= 6 {
					soPath := parts[5]
					// 获取库目录
					libDir := filepath.Dir(soPath)
					// 向上查找到 lib 目录
					for i := 0; i < 3; i++ {
						if filepath.Base(libDir) == "lib" {
							return libDir, nil
						}
						libDir = filepath.Dir(libDir)
					}
				}
			}
		}
	}

	return "", fmt.Errorf("无法找到 native library 目录")
}

// ExtractECSBinary 获取 ECS 二进制文件路径
// 在 Android 上，我们不需要"提取"，而是直接使用系统已安装的 native library
func ExtractECSBinary() (string, error) {
	// 获取 native library 目录
	libDir, err := findNativeLibraryDir()
	debugInfo := fmt.Sprintf("架构: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	if err != nil {
		debugInfo += fmt.Sprintf("查找 lib 目录失败: %v\n", err)
	} else {
		debugInfo += fmt.Sprintf("找到 lib 目录: %s\n", libDir)

		// 列出 lib 目录内容
		if entries, err := os.ReadDir(libDir); err == nil {
			debugInfo += fmt.Sprintf("lib 目录内容 (%d 项):\n", len(entries))
			for _, entry := range entries {
				entryType := "文件"
				if entry.IsDir() {
					entryType = "目录"
				}
				debugInfo += fmt.Sprintf("  - %s (%s)\n", entry.Name(), entryType)
			}
		}
	}

	// 库名称固定为 libgoecs.so
	libraryName := "libgoecs.so"

	// 尝试的子目录（Android ABI 名称）
	abiDirs := []string{
		"", // 直接在 lib 目录
	}

	// 根据架构添加 ABI 目录
	switch runtime.GOARCH {
	case "arm64":
		abiDirs = append(abiDirs, "arm64-v8a", "arm64")
	case "arm":
		abiDirs = append(abiDirs, "armeabi-v7a", "armeabi", "arm")
	case "amd64":
		abiDirs = append(abiDirs, "x86_64", "x86-64")
	case "386":
		abiDirs = append(abiDirs, "x86")
	}

	// 尝试所有可能的路径组合
	var checkedPaths []string

	if err == nil {
		for _, abiDir := range abiDirs {
			baseDir := libDir
			if abiDir != "" {
				baseDir = filepath.Join(libDir, abiDir)
			}

			ecsPath := filepath.Join(baseDir, libraryName)
			checkedPaths = append(checkedPaths, ecsPath)

			if info, err := os.Stat(ecsPath); err == nil && !info.IsDir() {
				// 找到文件，确保有执行权限
				if err := os.Chmod(ecsPath, 0755); err != nil {
					// 在某些 Android 版本上可能无法修改权限，但这通常不是问题
				}
				return ecsPath, nil
			}
		}
	}

	// 如果上述方法都失败，尝试在常见位置查找
	// 注意：不再查找 /system/lib，因为那是系统库位置
	fallbackPaths := []string{
		"/data/local/tmp/libgoecs.so", // 临时目录（需要 root）
	}

	for _, path := range fallbackPaths {
		checkedPaths = append(checkedPaths, path)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}

	// 未找到文件，返回详细错误信息
	recommendedABI := "arm64-v8a"
	if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
		recommendedABI = "x86_64"
	}

	return "", fmt.Errorf("找不到 ECS 二进制文件\n\n调试信息:\n%s\n已检查的路径:\n  %s\n\n请确保:\n1. ECS 二进制文件已编译为 Android 版本（Linux/%s）\n2. 文件已放置在 jniLibs/%s/libgoecs.so\n3. APK 已重新打包\n4. 当前架构: %s",
		debugInfo,
		strings.Join(checkedPaths, "\n  "),
		runtime.GOARCH,
		recommendedABI,
		runtime.GOARCH)
}

// CleanupECSBinary 清理函数
// 在 Android 上，native library 由系统管理，我们不需要清理
func CleanupECSBinary(path string) {
	// 不需要做任何事情
	// native library 由 Android 系统管理，应用卸载时会自动清理
}
