//go:build windows && arm64
// +build windows,arm64

package embedding

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed binaries/goecs-windows-arm64.exe
var ecsBinary []byte

func getECSBinary() ([]byte, error) {
	if len(ecsBinary) == 0 {
		return nil, fmt.Errorf("Windows ARM64 二进制文件未嵌入")
	}
	return ecsBinary, nil
}

func ExtractECSBinary() (string, error) {
	binary, err := getECSBinary()
	if err != nil {
		return "", err
	}

	tmpDir := os.TempDir()
	ecsPath := filepath.Join(tmpDir, "goecs.exe")

	if err := os.WriteFile(ecsPath, binary, 0755); err != nil {
		return "", fmt.Errorf("写入二进制文件失败: %v", err)
	}

	return ecsPath, nil
}

func CleanupECSBinary(path string) {
	if path != "" {
		os.Remove(path)
	}
}
