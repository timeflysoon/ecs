package main

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// onPresetChanged 预设模式改变时的处理
func (ui *TestUI) onPresetChanged(preset string) {
	switch preset {
	case "1. 融合怪完全体(能测全测)":
		ui.setAllChecks(true)
	case "2. 极简版(系统+CPU+内存+磁盘+5测速节点)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.speedCheck.Checked = true
	case "3. 精简版(系统+CPU+内存+磁盘+基础解锁+5测速节点)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.unlockCheck.Checked = true
		ui.nt3Check.Checked = true
		ui.speedCheck.Checked = true
	case "4. 精简网络版(系统+CPU+内存+磁盘+回程+路由+5测速节点)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.backtraceCheck.Checked = true
		ui.nt3Check.Checked = true
		ui.speedCheck.Checked = true
	case "5. 精简解锁版(系统+CPU+内存+磁盘IO+御三家+常用流媒体+5测速节点)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.commCheck.Checked = true
		ui.unlockCheck.Checked = true
		ui.speedCheck.Checked = true
	case "6. 仅网络测试(IP质量+5测速节点)":
		ui.setAllChecks(false)
		ui.securityCheck.Checked = true
		ui.speedCheck.Checked = true
		ui.backtraceCheck.Checked = true
		ui.nt3Check.Checked = true
		ui.pingCheck.Checked = true
	case "7. 仅解锁测试(基础解锁+常用流媒体解锁)":
		ui.setAllChecks(false)
		ui.commCheck.Checked = true
		ui.unlockCheck.Checked = true
	case "8. 仅硬件测试(系统+CPU+内存+dd磁盘+fio磁盘)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.diskMethodSelect.Selected = "fio"
	case "9. IP质量测试(IP测试+15数据库+邮件端口)":
		ui.setAllChecks(false)
		ui.securityCheck.Checked = true
		ui.emailCheck.Checked = true
	default: // 自定义
		return
	}
	ui.refreshAllChecks()
}

// setAllChecks 设置所有测试项的选中状态
func (ui *TestUI) setAllChecks(checked bool) {
	ui.basicCheck.Checked = checked
	ui.cpuCheck.Checked = checked
	ui.memoryCheck.Checked = checked
	ui.diskCheck.Checked = checked
	ui.commCheck.Checked = checked
	ui.unlockCheck.Checked = checked
	ui.securityCheck.Checked = checked
	ui.emailCheck.Checked = checked
	ui.backtraceCheck.Checked = checked
	ui.nt3Check.Checked = checked
	ui.speedCheck.Checked = checked
	ui.pingCheck.Checked = checked
	ui.refreshAllChecks()
}

// refreshAllChecks 刷新所有测试项的显示
func (ui *TestUI) refreshAllChecks() {
	ui.basicCheck.Refresh()
	ui.cpuCheck.Refresh()
	ui.memoryCheck.Refresh()
	ui.diskCheck.Refresh()
	ui.commCheck.Refresh()
	ui.unlockCheck.Refresh()
	ui.securityCheck.Refresh()
	ui.emailCheck.Refresh()
	ui.backtraceCheck.Refresh()
	ui.nt3Check.Refresh()
	ui.speedCheck.Refresh()
	ui.pingCheck.Refresh()
}

// startTests 开始执行测试
func (ui *TestUI) startTests() {
	ui.mu.Lock()
	if ui.isRunning {
		ui.mu.Unlock()
		return
	}
	ui.isRunning = true
	ui.mu.Unlock()

	if !ui.hasSelectedTests() {
		dialog.ShowInformation("提示", "请至少选择一项测试！", ui.window)
		ui.mu.Lock()
		ui.isRunning = false
		ui.mu.Unlock()
		return
	}

	ui.startButton.Disable()
	ui.stopButton.Enable()
	ui.clearButton.Disable()
	ui.progressBar.Show()
	ui.statusLabel.SetText("测试运行中...")
	// 清空结果显示
	ui.resultText.Segments = []widget.RichTextSegment{}
	ui.resultText.Refresh()

	ui.cancelCtx, ui.cancelFn = context.WithCancel(context.Background())
	go ui.runTests()
}

// stopTests 停止正在执行的测试
func (ui *TestUI) stopTests() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	if ui.cancelFn != nil {
		ui.cancelFn()
	}
	ui.statusLabel.SetText("测试已停止")
	ui.appendResult("\n\n========== 测试被用户中断 ==========\n")
	ui.resetUIState()
}

// clearResults 清空测试结果
func (ui *TestUI) clearResults() {
	ui.resultText.Segments = []widget.RichTextSegment{}
	ui.resultText.Refresh()
	ui.statusLabel.SetText("就绪")
	ui.progressBar.SetValue(0)
}

// exportResults 导出测试结果
func (ui *TestUI) exportResults() {
	// 从 RichText 提取纯文本内容
	var content string
	for _, seg := range ui.resultText.Segments {
		content += seg.Textual()
	}

	if content == "" {
		dialog.ShowInformation("提示", "没有可导出的结果", ui.window)
		return
	}

	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, ui.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		_, err = writer.Write([]byte(content))
		if err != nil {
			dialog.ShowError(err, ui.window)
			return
		}

		dialog.ShowInformation("成功", "结果已导出到: "+writer.URI().Path(), ui.window)
	}, ui.window)
}
