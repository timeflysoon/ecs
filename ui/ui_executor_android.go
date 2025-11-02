//go:build android

package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/oneclickvirt/ecs-android/embedding"
)

// createAndroidCommand 在 Android 平台创建命令
// 尝试使用 proot 运行 Linux 二进制文件，如果失败则直接运行
func (e *CommandExecutor) createAndroidCommand(args []string) (*exec.Cmd, string) {
	// 尝试获取 proot
	prootPath, err := embedding.GetProotPath()
	if err != nil {
		// 如果没有 proot，尝试直接运行（可能会失败）
		e.ui.Terminal.AppendText(fmt.Sprintf("⚠️  警告: 找不到 proot (%v)，尝试直接运行...\n\n", err))
		cmd := exec.CommandContext(e.cancelCtx, e.ecsPath, args...)
		cmdStr := fmt.Sprintf("执行命令: %s %s", e.ecsPath, strings.Join(args, " "))
		return cmd, cmdStr
	}

	// 使用 proot 运行
	e.ui.Terminal.AppendText("✓ 使用 proot 环境运行 Linux 二进制文件\n\n")

	prootCmd, prootArgs := embedding.BuildProotCommand(prootPath, e.ecsPath, args)
	cmd := exec.CommandContext(e.cancelCtx, prootCmd, prootArgs...)
	cmdStr := fmt.Sprintf("通过 proot 执行: %s %s", e.ecsPath, strings.Join(args, " "))

	// 设置环境变量
	if env, err := embedding.SetupProotEnvironment(); err == nil {
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return cmd, cmdStr
}
