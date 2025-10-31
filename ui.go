package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/oneclickvirt/CommonMediaTests/commediatests"
	"github.com/oneclickvirt/ecs/cputest"
	"github.com/oneclickvirt/ecs/disktest"
	"github.com/oneclickvirt/ecs/memorytest"
	"github.com/oneclickvirt/ecs/nexttrace"
	"github.com/oneclickvirt/ecs/speedtest"
	"github.com/oneclickvirt/ecs/unlocktest"
	"github.com/oneclickvirt/ecs/upstreams"
	"github.com/oneclickvirt/ecs/utils"
	"github.com/oneclickvirt/pingtest/pt"
	"github.com/oneclickvirt/portchecker/email"
)

const (
	// ecsVersion 使用核心包的版本号
	ecsVersion = "v0.1.91"
	width      = 80
)

type TestUI struct {
	app    fyne.App
	window fyne.Window

	// 测试选项复选框 - 完整支持所有测试项
	basicCheck     *widget.Check
	cpuCheck       *widget.Check
	memoryCheck    *widget.Check
	diskCheck      *widget.Check
	commCheck      *widget.Check // 御三家流媒体
	unlockCheck    *widget.Check // 跨国流媒体解锁
	securityCheck  *widget.Check // IP质量检测
	emailCheck     *widget.Check // 邮件端口检测
	backtraceCheck *widget.Check // 上游及回程线路
	nt3Check       *widget.Check // 三网回程路由
	speedCheck     *widget.Check // 网络测速
	pingCheck      *widget.Check // 三网PING值

	// 预设模式选择
	presetSelect *widget.Select

	// 配置选项
	languageSelect     *widget.Select
	cpuMethodSelect    *widget.Select
	memoryMethodSelect *widget.Select
	diskMethodSelect   *widget.Select
	diskPathEntry      *widget.Entry
	threadModeSelect   *widget.Select
	nt3LocationSelect  *widget.Select
	nt3TypeSelect      *widget.Select
	diskMultiCheck     *widget.Check
	spNumEntry         *widget.Entry

	// 控制按钮
	startButton *widget.Button
	stopButton  *widget.Button
	clearButton *widget.Button

	// 结果显示
	resultText  *widget.Entry
	progressBar *widget.ProgressBar
	statusLabel *widget.Label

	// 运行状态
	isRunning bool
	cancelCtx context.Context
	cancelFn  context.CancelFunc
	mu        sync.Mutex
}

func NewTestUI(app fyne.App) *TestUI {
	ui := &TestUI{
		app:    app,
		window: app.NewWindow("融合怪测试 - 完整版"),
	}

	// 设置窗口大小
	ui.window.Resize(fyne.NewSize(1200, 800))
	ui.window.SetPadded(true)

	ui.buildUI()
	return ui
}

func (ui *TestUI) buildUI() {
	// 创建选项卡
	tabs := container.NewAppTabs(
		container.NewTabItem("测试项目", ui.createTestOptionsTab()),
		container.NewTabItem("配置选项", ui.createConfigTab()),
		container.NewTabItem("测试结果", ui.createResultTab()),
	)

	ui.window.SetContent(tabs)
}

func (ui *TestUI) createTestOptionsTab() *fyne.Container {
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

	// 组合布局
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

	return container.NewHBox(
		container.NewScroll(leftColumn),
		container.NewScroll(rightColumn),
	)
}

func (ui *TestUI) createConfigTab() *fyne.Container {
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

	return container.NewHBox(
		container.NewScroll(leftColumn),
		container.NewScroll(rightColumn),
	)
}

func (ui *TestUI) createResultTab() *fyne.Container {
	// 状态标签
	ui.statusLabel = widget.NewLabel("就绪")

	// 进度条
	ui.progressBar = widget.NewProgressBar()
	ui.progressBar.Hide()

	// 结果文本框
	ui.resultText = widget.NewMultiLineEntry()
	ui.resultText.Wrapping = fyne.TextWrapWord
	ui.resultText.SetPlaceHolder("测试结果将显示在这里...")

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
	ui.resultText.SetText("")

	ui.cancelCtx, ui.cancelFn = context.WithCancel(context.Background())
	go ui.runTests()
}

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

func (ui *TestUI) clearResults() {
	ui.resultText.SetText("")
	ui.statusLabel.SetText("就绪")
	ui.progressBar.SetValue(0)
}

func (ui *TestUI) hasSelectedTests() bool {
	return ui.basicCheck.Checked ||
		ui.cpuCheck.Checked ||
		ui.memoryCheck.Checked ||
		ui.diskCheck.Checked ||
		ui.commCheck.Checked ||
		ui.unlockCheck.Checked ||
		ui.securityCheck.Checked ||
		ui.emailCheck.Checked ||
		ui.backtraceCheck.Checked ||
		ui.nt3Check.Checked ||
		ui.speedCheck.Checked ||
		ui.pingCheck.Checked
}

func (ui *TestUI) runTests() {
	defer func() {
		if r := recover(); r != nil {
			ui.appendResult(fmt.Sprintf("\n错误: %v\n", r))
		}
		ui.resetUIState()
	}()

	totalTests := ui.countSelectedTests()
	currentTest := 0

	language := "zh"
	if ui.languageSelect.Selected == "English" {
		language = "en"
	}

	startTime := time.Now()

	// 打印头部
	output := utils.PrintAndCapture(func() {
		utils.PrintHead(language, width, ecsVersion)
	}, "", "")
	ui.appendResult(output)

	// 网络检测
	preCheck := utils.CheckPublicAccess(3 * time.Second)

	// 执行各项测试
	var wg1, wg2, wg3 sync.WaitGroup
	var mediaInfo, emailInfo, ptInfo string

	// 基础信息测试
	if ui.basicCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "基础信息测试")
		ui.runBasicTest(language, preCheck)
	}

	// CPU 测试
	if ui.cpuCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "CPU 性能测试")
		ui.runCPUTest(language)
	}

	// 内存测试
	if ui.memoryCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "内存性能测试")
		ui.runMemoryTest(language)
	}

	// 磁盘测试
	if ui.diskCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "磁盘性能测试")
		ui.runDiskTest(language)
	}

	// 启动异步测试
	if ui.unlockCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			mediaInfo = unlocktest.MediaTest(language)
		}()
	}

	if ui.emailCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			emailInfo = email.EmailCheck()
		}()
	}

	if ui.pingCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			ptInfo = pt.PingTest()
		}()
	}

	// 御三家流媒体测试
	if ui.commCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "御三家流媒体测试")
		ui.runCommMediaTest(language)
	}

	// 跨国流媒体解锁测试
	if ui.unlockCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "跨国流媒体解锁测试")
		ui.runUnlockTest(language, &wg1, &mediaInfo)
	}

	// IP质量检测
	if ui.securityCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "IP质量检测")
		ui.runSecurityTest(language, preCheck)
	}

	// 邮件端口检测
	if ui.emailCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "邮件端口检测")
		ui.runEmailTest(language, &wg2, &emailInfo)
	}

	// 上游及回程线路检测
	if ui.backtraceCheck.Checked && preCheck.Connected && runtime.GOOS != "windows" && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "上游及回程线路检测")
		ui.runBacktraceTest(language)
	}

	// 三网回程路由检测
	if ui.nt3Check.Checked && preCheck.Connected && runtime.GOOS != "windows" && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "三网回程路由检测")
		ui.runNT3Test(language)
	}

	// 三网PING值检测
	if ui.pingCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "三网PING值检测")
		ui.runPingTest(language, &wg3, &ptInfo)
	}

	// 网络测速
	if ui.speedCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "网络测速")
		ui.runSpeedTest(language)
	}

	if !ui.isCancelled() {
		// 显示结束时间
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		minutes := int(duration.Minutes())
		seconds := int(duration.Seconds()) % 60
		currentTimeStr := endTime.Format("Mon Jan 2 15:04:05 MST 2006")

		timeOutput := utils.PrintAndCapture(func() {
			utils.PrintCenteredTitle("", width)
			if language == "zh" {
				fmt.Printf("花费          : %d 分 %d 秒\n", minutes, seconds)
				fmt.Printf("时间          : %s\n", currentTimeStr)
			} else {
				fmt.Printf("Cost    Time          : %d min %d sec\n", minutes, seconds)
				fmt.Printf("Current Time          : %s\n", currentTimeStr)
			}
			utils.PrintCenteredTitle("", width)
		}, "", "")

		ui.appendResult(timeOutput)
		ui.statusLabel.SetText("测试完成")
		dialog.ShowInformation("完成", "所有测试已完成！", ui.window)
	}
}

// 各测试函数实现
func (ui *TestUI) runBasicTest(language string, preCheck utils.NetCheckResult) {
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("系统基础信息", width)
		} else {
			utils.PrintCenteredTitle("System-Basic-Information", width)
		}

		var basicInfo, securityInfo string
		var nt3CheckType string = ui.nt3TypeSelect.Selected

		if preCheck.Connected && preCheck.StackType == "DualStack" {
			_, _, basicInfo, securityInfo, nt3CheckType = utils.BasicsAndSecurityCheck(language, nt3CheckType, ui.securityCheck.Checked)
		} else if preCheck.Connected && preCheck.StackType == "IPv4" {
			_, _, basicInfo, securityInfo, nt3CheckType = utils.BasicsAndSecurityCheck(language, "ipv4", ui.securityCheck.Checked)
		} else if preCheck.Connected && preCheck.StackType == "IPv6" {
			_, _, basicInfo, securityInfo, nt3CheckType = utils.BasicsAndSecurityCheck(language, "ipv6", ui.securityCheck.Checked)
		} else {
			_, _, basicInfo, securityInfo, nt3CheckType = utils.BasicsAndSecurityCheck(language, "", false)
		}

		fmt.Printf("%s", basicInfo)

		// 如果启用了安全检测但没有单独选中，这里也显示
		if ui.securityCheck.Checked && securityInfo != "" {
			fmt.Printf("%s", securityInfo)
		}
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runCPUTest(language string) {
	output := utils.PrintAndCapture(func() {
		realTestMethod, res := cputest.CpuTest(language, ui.cpuMethodSelect.Selected, ui.threadModeSelect.Selected)
		if language == "zh" {
			utils.PrintCenteredTitle(fmt.Sprintf("CPU测试-通过%s测试", realTestMethod), width)
		} else {
			utils.PrintCenteredTitle(fmt.Sprintf("CPU-Test--%s-Method", realTestMethod), width)
		}
		fmt.Print(res)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runMemoryTest(language string) {
	output := utils.PrintAndCapture(func() {
		realTestMethod, res := memorytest.MemoryTest(language, ui.memoryMethodSelect.Selected)
		if language == "zh" {
			utils.PrintCenteredTitle(fmt.Sprintf("内存测试-通过%s测试", realTestMethod), width)
		} else {
			utils.PrintCenteredTitle(fmt.Sprintf("Memory-Test--%s-Method", realTestMethod), width)
		}
		fmt.Print(res)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runDiskTest(language string) {
	output := utils.PrintAndCapture(func() {
		diskPath := ui.diskPathEntry.Text
		diskMethod := ui.diskMethodSelect.Selected
		diskMultiCheck := ui.diskMultiCheck.Checked
		autoChange := (diskMethod == "auto")

		realTestMethod, res := disktest.DiskTest(language, diskMethod, diskPath, diskMultiCheck, autoChange)
		if language == "zh" {
			utils.PrintCenteredTitle(fmt.Sprintf("硬盘测试-通过%s测试", realTestMethod), width)
		} else {
			utils.PrintCenteredTitle(fmt.Sprintf("Disk-Test--%s-Method", realTestMethod), width)
		}
		fmt.Print(res)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runCommMediaTest(language string) {
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("御三家流媒体解锁", width)
		} else {
			utils.PrintCenteredTitle("Common-Streaming-Media-Unlock", width)
		}
		fmt.Printf("%s", commediatests.MediaTests(language))
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runUnlockTest(language string, wg *sync.WaitGroup, mediaInfo *string) {
	wg.Wait()
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("跨国流媒体解锁", width)
		} else {
			utils.PrintCenteredTitle("Cross-Border-Streaming-Media-Unlock", width)
		}
		fmt.Printf("%s", *mediaInfo)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runSecurityTest(language string, preCheck utils.NetCheckResult) {
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("IP质量检测", width)
		} else {
			utils.PrintCenteredTitle("IP-Quality-Check", width)
		}

		var securityInfo string
		if preCheck.Connected {
			_, _, _, securityInfo, _ = utils.BasicsAndSecurityCheck(language, "", true)
		}
		fmt.Printf("%s", securityInfo)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runEmailTest(language string, wg *sync.WaitGroup, emailInfo *string) {
	wg.Wait()
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("邮件端口检测", width)
		} else {
			utils.PrintCenteredTitle("Email-Port-Check", width)
		}
		fmt.Println(*emailInfo)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runBacktraceTest(language string) {
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("上游及回程线路检测", width)
		} else {
			utils.PrintCenteredTitle("Upstream-and-Return-Route-Check", width)
		}
		upstreams.UpstreamsCheck()
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runNT3Test(language string) {
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("三网回程路由检测", width)
		} else {
			utils.PrintCenteredTitle("Three-Network-Return-Route-Check", width)
		}
		nexttrace.NextTrace3Check(language, ui.nt3LocationSelect.Selected, ui.nt3TypeSelect.Selected)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runPingTest(language string, wg *sync.WaitGroup, ptInfo *string) {
	wg.Wait()
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("三网ICMP的PING值检测", width)
		} else {
			utils.PrintCenteredTitle("Three-Network-ICMP-Ping-Check", width)
		}
		fmt.Println(*ptInfo)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) runSpeedTest(language string) {
	output := utils.PrintAndCapture(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("就近节点测速", width)
		} else {
			utils.PrintCenteredTitle("Speed-Test", width)
		}

		speedtest.ShowHead(language)
		speedtest.NearbySP()

		// 根据预设模式调整测速节点数
		spNum := 2
		if ui.spNumEntry.Text != "" {
			fmt.Sscanf(ui.spNumEntry.Text, "%d", &spNum)
		}

		speedtest.CustomSP("net", "global", 2, language)
		speedtest.CustomSP("net", "cu", spNum, language)
		speedtest.CustomSP("net", "ct", spNum, language)
		speedtest.CustomSP("net", "cmcc", spNum, language)
	}, "", "")
	ui.appendResult(output)
}

func (ui *TestUI) countSelectedTests() int {
	count := 0
	if ui.basicCheck.Checked {
		count++
	}
	if ui.cpuCheck.Checked {
		count++
	}
	if ui.memoryCheck.Checked {
		count++
	}
	if ui.diskCheck.Checked {
		count++
	}
	if ui.commCheck.Checked {
		count++
	}
	if ui.unlockCheck.Checked {
		count++
	}
	if ui.securityCheck.Checked {
		count++
	}
	if ui.emailCheck.Checked {
		count++
	}
	if ui.backtraceCheck.Checked {
		count++
	}
	if ui.nt3Check.Checked {
		count++
	}
	if ui.speedCheck.Checked {
		count++
	}
	if ui.pingCheck.Checked {
		count++
	}
	return count
}

func (ui *TestUI) updateProgress(current, total int, testName string) {
	progress := float64(current) / float64(total)
	ui.progressBar.SetValue(progress)
	ui.statusLabel.SetText(fmt.Sprintf("[%d/%d] %s", current, total, testName))
}

func (ui *TestUI) appendResult(text string) {
	currentText := ui.resultText.Text
	ui.resultText.SetText(currentText + text)
	ui.resultText.CursorRow = len(strings.Split(ui.resultText.Text, "\n"))
}

func (ui *TestUI) isCancelled() bool {
	select {
	case <-ui.cancelCtx.Done():
		return true
	default:
		return false
	}
}

func (ui *TestUI) resetUIState() {
	ui.mu.Lock()
	ui.isRunning = false
	ui.mu.Unlock()

	ui.startButton.Enable()
	ui.stopButton.Disable()
	ui.clearButton.Enable()
	ui.progressBar.Hide()
	ui.progressBar.SetValue(0)
}

func (ui *TestUI) exportResults() {
	content := ui.resultText.Text
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
