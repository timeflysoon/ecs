package ui

// hasSelectedTests 检查是否有选中的测试项
func (ui *TestUI) hasSelectedTests() bool {
	return ui.BasicCheck.Checked ||
		ui.CpuCheck.Checked ||
		ui.MemoryCheck.Checked ||
		ui.DiskCheck.Checked ||
		ui.CommCheck.Checked ||
		ui.UnlockCheck.Checked ||
		ui.SecurityCheck.Checked ||
		ui.EmailCheck.Checked ||
		ui.BacktraceCheck.Checked ||
		ui.Nt3Check.Checked ||
		ui.SpeedCheck.Checked ||
		ui.PingCheck.Checked
}

// isCancelled 检查测试是否被取消
func (ui *TestUI) isCancelled() bool {
	select {
	case <-ui.CancelCtx.Done():
		return true
	default:
		return false
	}
}

// resetUIState 重置UI状态
func (ui *TestUI) resetUIState() {
	ui.Mu.Lock()
	ui.IsRunning = false
	ui.Mu.Unlock()

	ui.StartButton.Enable()
	ui.StopButton.Disable()
	ui.ProgressBar.Hide()
	ui.ProgressBar.SetValue(0)
}
