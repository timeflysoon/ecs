package ui

import (
	"fmt"
	"time"
)

// runTestsWithExecutor 使用命令执行器运行测试
func (ui *TestUI) runTestsWithExecutor() {
	defer func() {
		if r := recover(); r != nil {
			ui.Terminal.AppendText(fmt.Sprintf("\n错误: %v\n", r))
		}
		ui.resetUIState()
	}()

	startTime := time.Now()

	// 创建命令执行器
	executor, err := NewCommandExecutor(ui, ui.CancelCtx)
	if err != nil {
		ui.Terminal.AppendText(fmt.Sprintf("错误: %v\n", err))
		ui.StatusLabel.SetText("测试失败")
		return
	}
	defer executor.Cleanup()

	// 显示将要执行的命令
	cmdPreview := executor.GetCommandPreview()
	ui.Terminal.AppendText(fmt.Sprintf("执行命令: %s\n", cmdPreview))

	// 更新进度
	ui.ProgressBar.SetValue(0.1)
	ui.StatusLabel.SetText("正在执行测试...")

	// 执行测试（输出会实时显示在terminal widget中）
	err = executor.Execute()

	// 显示结束信息
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	_ = duration // 避免未使用警告

	if err != nil {
		ui.Terminal.AppendText(fmt.Sprintf("\n错误: %v\n", err))
		ui.StatusLabel.SetText("测试失败")
	} else if ui.isCancelled() {
		ui.Terminal.AppendText("\n测试被用户中断\n")
		ui.StatusLabel.SetText("测试已停止")
	} else {
		ui.StatusLabel.SetText("测试完成")
		ui.ProgressBar.SetValue(1.0)

		// 如果启用了日志，自动刷新日志内容
		if ui.LogCheck != nil && ui.LogCheck.Checked {
			time.Sleep(500 * time.Millisecond) // 等待日志文件写入完成
			ui.refreshLogFromFile()
		}
	}

	ui.Terminal.AppendText("==========================================\n")
}
