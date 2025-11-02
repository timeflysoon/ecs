package ui

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/oneclickvirt/ecs-android/embedding"
)

// CommandExecutor 负责执行嵌入的 ecs 二进制文件
type CommandExecutor struct {
	ecsPath   string
	ui        *TestUI
	mu        sync.Mutex
	cancelCtx context.Context
}

// NewCommandExecutor 创建新的命令执行器
func NewCommandExecutor(testUI *TestUI, ctx context.Context) (*CommandExecutor, error) {
	// 提取嵌入的二进制文件
	ecsPath, err := embedding.ExtractECSBinary()
	if err != nil {
		return nil, fmt.Errorf("提取 ECS 二进制文件失败: %v", err)
	}

	return &CommandExecutor{
		ecsPath:   ecsPath,
		ui:        testUI,
		cancelCtx: ctx,
	}, nil
}

// Cleanup 清理临时文件
func (e *CommandExecutor) Cleanup() {
	if e.ecsPath != "" {
		embedding.CleanupECSBinary(e.ecsPath)
	}
}

// buildCommandArgs 根据UI选项构建命令参数
func (e *CommandExecutor) buildCommandArgs() []string {
	args := []string{}

	// 禁用菜单模式（GUI模式必须禁用CLI菜单）
	args = append(args, "-menu=false")

	// 语言参数
	if e.ui.LanguageSelect.Selected == "English" {
		args = append(args, "-l", "en")
	} else {
		args = append(args, "-l", "zh")
	}

	// 基础测试
	if !e.ui.BasicCheck.Checked {
		args = append(args, "-basic=false")
	}

	// CPU 测试
	if !e.ui.CpuCheck.Checked {
		args = append(args, "-cpu=false")
	} else {
		// CPU测试方法
		args = append(args, "-cpum", e.ui.CpuMethodSelect.Selected)
		// CPU线程模式
		args = append(args, "-cput", e.ui.ThreadModeSelect.Selected)
	}

	// 内存测试
	if !e.ui.MemoryCheck.Checked {
		args = append(args, "-memory=false")
	} else {
		// 内存测试方法
		args = append(args, "-memorym", e.ui.MemoryMethodSelect.Selected)
	}

	// 磁盘测试
	if !e.ui.DiskCheck.Checked {
		args = append(args, "-disk=false")
	} else {
		// 磁盘测试方法
		args = append(args, "-diskm", e.ui.DiskMethodSelect.Selected)
		// 磁盘测试路径
		if e.ui.DiskPathEntry.Text != "" {
			args = append(args, "-diskp", e.ui.DiskPathEntry.Text)
		}
		// 多磁盘检查
		if e.ui.DiskMultiCheck.Checked {
			args = append(args, "-diskmc=true")
		} else {
			args = append(args, "-diskmc=false")
		}
	}

	// 御三家流媒体测试
	if !e.ui.CommCheck.Checked {
		args = append(args, "-comm=false")
	}

	// 跨国流媒体解锁测试
	if !e.ui.UnlockCheck.Checked {
		args = append(args, "-ut=false")
	}

	// IP质量检测
	if !e.ui.SecurityCheck.Checked {
		args = append(args, "-security=false")
	}

	// 邮件端口检测
	if !e.ui.EmailCheck.Checked {
		args = append(args, "-email=false")
	}

	// 上游及回程线路检测
	if !e.ui.BacktraceCheck.Checked {
		args = append(args, "-backtrace=false")
	}

	// 三网回程路由检测
	if !e.ui.Nt3Check.Checked {
		args = append(args, "-nt3=false")
	} else {
		// NT3测试位置
		args = append(args, "-nt3loc", e.ui.Nt3LocationSelect.Selected)
		// NT3测试类型
		args = append(args, "-nt3t", e.ui.Nt3TypeSelect.Selected)
	}

	// 网络测速
	if !e.ui.SpeedCheck.Checked {
		args = append(args, "-speed=false")
	} else {
		// 每个运营商的服务器数量
		if e.ui.SpNumEntry.Text != "" {
			args = append(args, "-spnum", e.ui.SpNumEntry.Text)
		}
	}

	// 上传结果（默认为true，GUI中可能需要控制）
	// 如果未来添加上传选项，可以在这里处理
	// args = append(args, "-upload=true")

	// 启用日志
	if e.ui.LogCheck != nil && e.ui.LogCheck.Checked {
		args = append(args, "-log")
	}

	return args
}

// Execute 执行命令（输出直接到终端widget）
func (e *CommandExecutor) Execute() error {
	e.mu.Lock()
	args := e.buildCommandArgs()
	e.mu.Unlock()

	var cmd *exec.Cmd
	var cmdStr string

	// 在 Android 上使用 proot 运行 Linux 二进制文件
	if runtime.GOOS == "android" {
		// Android 平台：尝试使用 proot，如果失败则直接运行
		cmd, cmdStr = e.createAndroidCommand(args)
	} else {
		// 非 Android 平台直接运行
		cmd = exec.CommandContext(e.cancelCtx, e.ecsPath, args...)
		cmdStr = fmt.Sprintf("执行命令: %s %s", e.ecsPath, strings.Join(args, " "))
	}

	e.ui.Terminal.AppendText(cmdStr + "\n\n")

	// 直接将stdout和stderr连接到终端widget
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建stdout管道失败: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建stderr管道失败: %v", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}

	// 实时读取并显示输出
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		buf := make([]byte, 4096)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				e.ui.Terminal.AppendText(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		buf := make([]byte, 4096)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				e.ui.Terminal.AppendText(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	// 等待输出读取完成
	wg.Wait()

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		if e.cancelCtx.Err() == context.Canceled {
			return fmt.Errorf("测试已被取消")
		}
		return fmt.Errorf("命令执行失败: %v", err)
	}

	return nil
}

// GetCommandPreview 获取将要执行的完整命令（用于调试）
func (e *CommandExecutor) GetCommandPreview() string {
	args := e.buildCommandArgs()
	return e.ecsPath + " " + strings.Join(args, " ")
}
