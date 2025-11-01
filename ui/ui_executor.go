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

	// 清空终端并显示开始信息
	ui.Terminal.AppendText("==========================================\n")
	ui.Terminal.AppendText("  融合怪测试 - 开始执行\n")
	ui.Terminal.AppendText("==========================================\n\n")

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
	ui.Terminal.AppendText(fmt.Sprintf("执行命令: %s\n\n", cmdPreview))
	ui.Terminal.AppendText("==========================================\n\n")

	// 更新进度
	ui.ProgressBar.SetValue(0.1)
	ui.StatusLabel.SetText("正在执行测试...")

	// 执行测试（输出会实时显示在terminal widget中）
	err = executor.Execute()

	// 显示结束信息
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60

	ui.Terminal.AppendText("\n\n==========================================\n")

	if err != nil {
		ui.Terminal.AppendText(fmt.Sprintf("错误: %v\n", err))
		ui.StatusLabel.SetText("测试失败")
	} else if ui.isCancelled() {
		ui.Terminal.AppendText("测试被用户中断\n")
		ui.StatusLabel.SetText("测试已停止")
	} else {
		language := "zh"
		if ui.LanguageSelect.Selected == "English" {
			language = "en"
		}

		if language == "zh" {
			ui.Terminal.AppendText(fmt.Sprintf("花费时间: %d 分 %d 秒\n", minutes, seconds))
		} else {
			ui.Terminal.AppendText(fmt.Sprintf("Cost Time: %d min %d sec\n", minutes, seconds))
		}

		ui.StatusLabel.SetText("测试完成")
		ui.ProgressBar.SetValue(1.0)

		// 如果启用了日志，自动刷新日志内容
		if ui.LogCheck != nil && ui.LogCheck.Checked {
			time.Sleep(500 * time.Millisecond) // 等待日志文件写入完成
			ui.refreshLog()
		}
	}

	ui.Terminal.AppendText("==========================================\n")
}
