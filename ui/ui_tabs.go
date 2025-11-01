package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// createOptionsPanel 创建选项面板（测试项目 + 配置选项整合在一起）
func (ui *TestUI) createOptionsPanel() fyne.CanvasObject {
	// 预设模式选择
	ui.PresetSelect = widget.NewSelect(
		[]string{
			"自定义",
			"1. 融合怪完全体(能测全测)",
			"2. 极简版(系统+CPU+内存+磁盘+5测速节点)",
			"3. 精简版(系统+CPU+内存+磁盘+基础解锁+5测速节点)",
			"4. 精简网络版(系统+CPU+内存+磁盘+回程+路由+5测速节点)",
			"5. 精简解锁版(系统+CPU+内存+磁盘IO+御三家+常用流媒体+5测速节点)",
			"6. 仅网络测试(IP质量+5测速节点)",
			"7. 仅解锁测试(基础解锁+常用流媒体解锁)",
			"8. 仅硬件测试(系统+CPU+内存+dd磁盘+fio磁盘)",
			"9. IP质量测试(IP测试+15数据库+邮件端口)",
		},
		ui.onPresetChanged,
	)
	ui.PresetSelect.Selected = "自定义"

	presetSection := widget.NewCard("预设模式", "快速选择测试组合", ui.PresetSelect)

	// === 测试项目复选框 ===
	ui.BasicCheck = widget.NewCheck("基础信息测试", nil)
	ui.BasicCheck.Checked = true

	ui.CpuCheck = widget.NewCheck("CPU 性能测试", nil)
	ui.CpuCheck.Checked = true

	ui.MemoryCheck = widget.NewCheck("内存性能测试", nil)
	ui.MemoryCheck.Checked = true

	ui.DiskCheck = widget.NewCheck("磁盘性能测试", nil)
	ui.DiskCheck.Checked = true

	ui.CommCheck = widget.NewCheck("御三家流媒体测试", nil)
	ui.CommCheck.Checked = false

	ui.UnlockCheck = widget.NewCheck("跨国流媒体解锁测试", nil)
	ui.UnlockCheck.Checked = false

	ui.SecurityCheck = widget.NewCheck("IP质量检测", nil)
	ui.SecurityCheck.Checked = false

	ui.EmailCheck = widget.NewCheck("邮件端口检测", nil)
	ui.EmailCheck.Checked = false

	ui.BacktraceCheck = widget.NewCheck("上游及回程线路检测", nil)
	ui.BacktraceCheck.Checked = false

	ui.Nt3Check = widget.NewCheck("三网回程路由检测", nil)
	ui.Nt3Check.Checked = false

	ui.SpeedCheck = widget.NewCheck("网络测速", nil)
	ui.SpeedCheck.Checked = false

	ui.PingCheck = widget.NewCheck("三网PING值检测", nil)
	ui.PingCheck.Checked = false

	ui.LogCheck = widget.NewCheck("启用日志记录", nil)
	ui.LogCheck.Checked = false

	// 全选/取消全选按钮
	selectAllBtn := widget.NewButton("全选", func() {
		ui.setAllChecks(true)
	})

	deselectAllBtn := widget.NewButton("取消全选", func() {
		ui.setAllChecks(false)
	})

	buttonRow := container.NewHBox(selectAllBtn, deselectAllBtn)

	// 测试项目分组 - 使用网格布局，每行2个
	basicTests := container.NewVBox(
		ui.BasicCheck,
		ui.CpuCheck,
		ui.MemoryCheck,
		ui.DiskCheck,
	)

	networkTests := container.NewVBox(
		ui.SpeedCheck,
		ui.SecurityCheck,
		ui.EmailCheck,
		ui.BacktraceCheck,
	)

	advancedTests := container.NewVBox(
		ui.Nt3Check,
		ui.PingCheck,
		ui.CommCheck,
		ui.UnlockCheck,
	)

	testsGrid := container.NewGridWithColumns(3,
		basicTests,
		networkTests,
		advancedTests,
	)

	testsSection := widget.NewCard("测试项目", "", container.NewVBox(
		buttonRow,
		testsGrid,
	))

	// === 配置选项 ===
	configSection := ui.createConfigSection()

	// 整合所有内容
	allContent := container.NewVBox(
		presetSection,
		testsSection,
		configSection,
	)

	return allContent
}

// createConfigSection 创建配置选项区域
func (ui *TestUI) createConfigSection() fyne.CanvasObject {
	// 语言选择
	ui.LanguageSelect = widget.NewSelect(
		[]string{"中文", "English"},
		func(value string) {},
	)
	ui.LanguageSelect.Selected = "中文"

	// CPU 配置
	ui.CpuMethodSelect = widget.NewSelect(
		[]string{"sysbench", "geekbench", "winsat"},
		func(value string) {},
	)
	ui.CpuMethodSelect.Selected = "sysbench"

	ui.ThreadModeSelect = widget.NewSelect(
		[]string{"single", "multi"},
		func(value string) {},
	)
	ui.ThreadModeSelect.Selected = "multi"

	// 内存配置
	ui.MemoryMethodSelect = widget.NewSelect(
		[]string{"auto", "stream", "sysbench", "dd", "winsat"},
		func(value string) {},
	)
	ui.MemoryMethodSelect.Selected = "auto"

	// 磁盘配置
	ui.DiskMethodSelect = widget.NewSelect(
		[]string{"auto", "fio", "dd", "winsat"},
		func(value string) {},
	)
	ui.DiskMethodSelect.Selected = "auto"

	ui.DiskPathEntry = widget.NewEntry()
	ui.DiskPathEntry.SetPlaceHolder("/tmp 或留空自动检测")

	ui.DiskMultiCheck = widget.NewCheck("启用多磁盘检测", nil)
	ui.DiskMultiCheck.Checked = false

	// NT3 配置
	ui.Nt3LocationSelect = widget.NewSelect(
		[]string{"GZ", "SH", "BJ", "CD", "ALL"},
		func(value string) {},
	)
	ui.Nt3LocationSelect.Selected = "GZ"

	ui.Nt3TypeSelect = widget.NewSelect(
		[]string{"ipv4", "ipv6", "both"},
		func(value string) {},
	)
	ui.Nt3TypeSelect.Selected = "ipv4"

	// 测速配置
	ui.SpNumEntry = widget.NewEntry()
	ui.SpNumEntry.SetText("2")
	ui.SpNumEntry.SetPlaceHolder("每运营商测速节点数")

	// 使用表单布局更紧凑
	configForm := container.NewVBox(
		widget.NewLabel("通用配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("语言:"),
			ui.LanguageSelect,
		),
		ui.LogCheck, // 日志选项
		widget.NewSeparator(),
		widget.NewLabel("CPU配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("测试方法:"),
			ui.CpuMethodSelect,
			widget.NewLabel("线程模式:"),
			ui.ThreadModeSelect,
		),
		widget.NewSeparator(),
		widget.NewLabel("内存配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("测试方法:"),
			ui.MemoryMethodSelect,
		),
		widget.NewSeparator(),
		widget.NewLabel("磁盘配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("测试方法:"),
			ui.DiskMethodSelect,
			widget.NewLabel("测试路径:"),
			ui.DiskPathEntry,
		),
		ui.DiskMultiCheck,
		widget.NewSeparator(),
		widget.NewLabel("三网回程配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("测试地点:"),
			ui.Nt3LocationSelect,
			widget.NewLabel("测试类型:"),
			ui.Nt3TypeSelect,
		),
		widget.NewSeparator(),
		widget.NewLabel("测速配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("节点数/运营商:"),
			ui.SpNumEntry,
		),
	)

	return widget.NewCard("详细配置", "调整测试参数", configForm)
}
