package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// createTestOptionsTab 创建测试项目选择选项卡
func (ui *TestUI) createTestOptionsTab() fyne.CanvasObject {
	// 预设模式选择
	ui.presetSelect = widget.NewSelect(
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
	ui.presetSelect.Selected = "自定义"

	presetCard := widget.NewCard("预设模式", "选择预设的测试组合", ui.presetSelect)

	// 创建所有测试项复选框
	ui.basicCheck = widget.NewCheck("基础信息测试", nil)
	ui.basicCheck.Checked = true

	ui.cpuCheck = widget.NewCheck("CPU 性能测试", nil)
	ui.cpuCheck.Checked = true

	ui.memoryCheck = widget.NewCheck("内存性能测试", nil)
	ui.memoryCheck.Checked = true

	ui.diskCheck = widget.NewCheck("磁盘性能测试", nil)
	ui.diskCheck.Checked = true

	ui.commCheck = widget.NewCheck("御三家流媒体测试", nil)
	ui.commCheck.Checked = false

	ui.unlockCheck = widget.NewCheck("跨国流媒体解锁测试", nil)
	ui.unlockCheck.Checked = false

	ui.securityCheck = widget.NewCheck("IP质量检测", nil)
	ui.securityCheck.Checked = false

	ui.emailCheck = widget.NewCheck("邮件端口检测", nil)
	ui.emailCheck.Checked = false

	ui.backtraceCheck = widget.NewCheck("上游及回程线路检测", nil)
	ui.backtraceCheck.Checked = false

	ui.nt3Check = widget.NewCheck("三网回程路由检测", nil)
	ui.nt3Check.Checked = false

	ui.speedCheck = widget.NewCheck("网络测速", nil)
	ui.speedCheck.Checked = false

	ui.pingCheck = widget.NewCheck("三网PING值检测", nil)
	ui.pingCheck.Checked = false

	// 全选/取消全选按钮
	selectAllBtn := widget.NewButton("全选", func() {
		ui.setAllChecks(true)
	})

	deselectAllBtn := widget.NewButton("取消全选", func() {
		ui.setAllChecks(false)
	})

	buttonRow := container.NewHBox(selectAllBtn, deselectAllBtn)

	// 分组显示
	basicGroup := widget.NewCard("基础测试", "", container.NewVBox(
		ui.basicCheck,
		ui.cpuCheck,
		ui.memoryCheck,
		ui.diskCheck,
	))

	networkGroup := widget.NewCard("网络测试", "", container.NewVBox(
		ui.speedCheck,
		ui.securityCheck,
		ui.emailCheck,
		ui.backtraceCheck,
		ui.nt3Check,
		ui.pingCheck,
	))

	mediaGroup := widget.NewCard("流媒体测试", "", container.NewVBox(
		ui.commCheck,
		ui.unlockCheck,
	))

	// 控制按钮
	ui.startButton = widget.NewButton("开始测试", ui.startTests)
	ui.startButton.Importance = widget.HighImportance

	ui.stopButton = widget.NewButton("停止测试", ui.stopTests)
	ui.stopButton.Disable()

	ui.clearButton = widget.NewButton("清空结果", ui.clearResults)

	controlButtons := container.NewVBox(
		ui.startButton,
		container.NewGridWithColumns(2, ui.stopButton, ui.clearButton),
	)

	controlCard := widget.NewCard("控制", "", controlButtons)

	// 组合布局 - 使用 Split 容器确保左右平均分配空间
	leftColumn := container.NewVBox(
		presetCard,
		buttonRow,
		basicGroup,
	)

	rightColumn := container.NewVBox(
		networkGroup,
		mediaGroup,
		controlCard,
	)

	leftScroll := container.NewScroll(leftColumn)
	leftScroll.SetMinSize(fyne.NewSize(400, 0))

	rightScroll := container.NewScroll(rightColumn)
	rightScroll.SetMinSize(fyne.NewSize(400, 0))

	return container.NewHSplit(leftScroll, rightScroll)
}

// createConfigTab 创建配置选项选项卡
func (ui *TestUI) createConfigTab() fyne.CanvasObject {
	// 语言选择
	ui.languageSelect = widget.NewSelect(
		[]string{"中文", "English"},
		func(value string) {},
	)
	ui.languageSelect.Selected = "中文"

	// CPU 配置
	ui.cpuMethodSelect = widget.NewSelect(
		[]string{"sysbench", "geekbench", "winsat"},
		func(value string) {},
	)
	ui.cpuMethodSelect.Selected = "sysbench"

	ui.threadModeSelect = widget.NewSelect(
		[]string{"single", "multi"},
		func(value string) {},
	)
	ui.threadModeSelect.Selected = "multi"

	cpuCard := widget.NewCard("CPU测试配置", "", container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("测试方法:"), nil, ui.cpuMethodSelect),
		container.NewBorder(nil, nil, widget.NewLabel("线程模式:"), nil, ui.threadModeSelect),
	))

	// 内存配置
	ui.memoryMethodSelect = widget.NewSelect(
		[]string{"auto", "stream", "sysbench", "dd", "winsat"},
		func(value string) {},
	)
	ui.memoryMethodSelect.Selected = "auto"

	memoryCard := widget.NewCard("内存测试配置", "", container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("测试方法:"), nil, ui.memoryMethodSelect),
	))

	// 磁盘配置
	ui.diskMethodSelect = widget.NewSelect(
		[]string{"auto", "fio", "dd", "winsat"},
		func(value string) {},
	)
	ui.diskMethodSelect.Selected = "auto"

	ui.diskPathEntry = widget.NewEntry()
	ui.diskPathEntry.SetPlaceHolder("/tmp 或留空自动检测")

	ui.diskMultiCheck = widget.NewCheck("启用多磁盘检测", nil)
	ui.diskMultiCheck.Checked = false

	diskCard := widget.NewCard("磁盘测试配置", "", container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("测试方法:"), nil, ui.diskMethodSelect),
		container.NewBorder(nil, nil, widget.NewLabel("测试路径:"), nil, ui.diskPathEntry),
		ui.diskMultiCheck,
	))

	// NT3 配置
	ui.nt3LocationSelect = widget.NewSelect(
		[]string{"GZ", "SH", "BJ", "CD", "ALL"},
		func(value string) {},
	)
	ui.nt3LocationSelect.Selected = "GZ"

	ui.nt3TypeSelect = widget.NewSelect(
		[]string{"ipv4", "ipv6", "both"},
		func(value string) {},
	)
	ui.nt3TypeSelect.Selected = "ipv4"

	nt3Card := widget.NewCard("三网回程配置", "", container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("测试地点:"), nil, ui.nt3LocationSelect),
		container.NewBorder(nil, nil, widget.NewLabel("测试类型:"), nil, ui.nt3TypeSelect),
	))

	// 测速配置
	ui.spNumEntry = widget.NewEntry()
	ui.spNumEntry.SetText("2")
	ui.spNumEntry.SetPlaceHolder("每运营商测速节点数")

	speedCard := widget.NewCard("测速配置", "", container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("节点数/运营商:"), nil, ui.spNumEntry),
	))

	// 通用配置
	generalCard := widget.NewCard("通用配置", "", container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("语言:"), nil, ui.languageSelect),
	))

	leftColumn := container.NewVBox(
		generalCard,
		cpuCard,
		memoryCard,
	)

	rightColumn := container.NewVBox(
		diskCard,
		nt3Card,
		speedCard,
	)

	leftScroll := container.NewScroll(leftColumn)
	leftScroll.SetMinSize(fyne.NewSize(400, 0))

	rightScroll := container.NewScroll(rightColumn)
	rightScroll.SetMinSize(fyne.NewSize(400, 0))

	return container.NewHSplit(leftScroll, rightScroll)
}

// createResultTab 创建测试结果显示选项卡
func (ui *TestUI) createResultTab() fyne.CanvasObject {
	// 状态标签
	ui.statusLabel = widget.NewLabel("就绪")

	// 进度条
	ui.progressBar = widget.NewProgressBar()
	ui.progressBar.Hide()

	// 结果文本 - 使用 RichText 支持富文本和颜色
	ui.resultText = widget.NewRichText()
	ui.resultText.Wrapping = fyne.TextWrapOff

	// 导出按钮
	exportButton := widget.NewButton("导出结果", ui.exportResults)

	topBar := container.NewBorder(
		nil, nil,
		ui.statusLabel,
		exportButton,
		ui.progressBar,
	)

	return container.NewBorder(
		topBar,
		nil, nil, nil,
		container.NewScroll(ui.resultText),
	)
}
