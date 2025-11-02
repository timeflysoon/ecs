//go:build !android

package ui

import (
	"fmt"
	"os/exec"
	"strings"
)

// createAndroidCommand 在非 Android 平台的占位函数
// 这个函数在非 Android 平台不会被调用，但需要存在以满足编译
func (e *CommandExecutor) createAndroidCommand(args []string) (*exec.Cmd, string) {
	// 非 Android 平台直接运行
	cmd := exec.CommandContext(e.cancelCtx, e.ecsPath, args...)
	cmdStr := fmt.Sprintf("执行命令: %s %s", e.ecsPath, strings.Join(args, " "))
	return cmd, cmdStr
}
