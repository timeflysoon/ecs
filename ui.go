package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/oneclickvirt/ecs/cputest"
	"github.com/oneclickvirt/ecs/disktest"
	"github.com/oneclickvirt/ecs/memorytest"
	"github.com/oneclickvirt/ecs/speedtest"
	"github.com/oneclickvirt/ecs/unlocktest"
	"github.com/oneclickvirt/ecs/utils"
)

type TestUI struct {
	app    fyne.App
	window fyne.Window

	// 测试选项复选框
	basicCheck  *widget.Check
	cpuCheck    *widget.Check
	memoryCheck *widget.Check
	diskCheck   *widget.Check
	speedCheck  *widget.Check
	unlockCheck *widget.Check
	routeCheck  *widget.Check

	// 配置选项
	languageSelect   *widget.Select
	cpuMethodSelect  *widget.Select
	diskPathEntry    *widget.Entry
	threadModeSelect *widget.Select

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
		window: app.NewWindow("融合怪测试"),
	}

	// 设置窗口大小 - 支持桌面和移动设备
	// 移动设备会自动全屏
	ui.window.Resize(fyne.NewSize(900, 700))
	ui.window.SetPadded(true)

	ui.buildUI()
	return ui
}

func (ui *TestUI) buildUI() {
	// 创建测试选项区域
	testOptionsGroup := ui.createTestOptions()

	// 创建配置选项区域
	configGroup := ui.createConfigOptions()

	// 创建控制按钮区域
	controlButtons := ui.createControlButtons()

	// 创建结果显示区域
	resultArea := ui.createResultArea()

	// 左侧面板：选项和配置（添加内边距以适应移动设备）
	leftPanel := container.NewVBox(
		testOptionsGroup,
		widget.NewSeparator(),
		configGroup,
		widget.NewSeparator(),
		controlButtons,
	)

	// 右侧面板：结果显示
	rightPanel := resultArea

	// 使用分割容器 - 在移动设备上会自动调整为垂直布局
	split := container.NewHSplit(
		container.NewScroll(leftPanel),
		rightPanel,
	)
	split.Offset = 0.4 // 左侧占40%，为移动设备优化

	ui.window.SetContent(split)
}

func (ui *TestUI) createTestOptions() *widget.Card {
	ui.basicCheck = widget.NewCheck("基础信息测试", nil)
	ui.basicCheck.Checked = true

	ui.cpuCheck = widget.NewCheck("CPU 性能测试", nil)
	ui.cpuCheck.Checked = true

	ui.memoryCheck = widget.NewCheck("内存性能测试", nil)
	ui.memoryCheck.Checked = true

	ui.diskCheck = widget.NewCheck("磁盘性能测试", nil)
	ui.diskCheck.Checked = true

	ui.speedCheck = widget.NewCheck("网络测速", nil)
	ui.speedCheck.Checked = false

	ui.unlockCheck = widget.NewCheck("流媒体解锁测试", nil)
	ui.unlockCheck.Checked = false

	ui.routeCheck = widget.NewCheck("路由追踪测试", nil)
	ui.routeCheck.Checked = false

	// 全选/取消全选按钮
	selectAllBtn := widget.NewButton("全选", func() {
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.speedCheck.Checked = true
		ui.unlockCheck.Checked = true
		ui.routeCheck.Checked = true
		ui.basicCheck.Refresh()
		ui.cpuCheck.Refresh()
		ui.memoryCheck.Refresh()
		ui.diskCheck.Refresh()
		ui.speedCheck.Refresh()
		ui.unlockCheck.Refresh()
		ui.routeCheck.Refresh()
	})

	deselectAllBtn := widget.NewButton("取消全选", func() {
		ui.basicCheck.Checked = false
		ui.cpuCheck.Checked = false
		ui.memoryCheck.Checked = false
		ui.diskCheck.Checked = false
		ui.speedCheck.Checked = false
		ui.unlockCheck.Checked = false
		ui.routeCheck.Checked = false
		ui.basicCheck.Refresh()
		ui.cpuCheck.Refresh()
		ui.memoryCheck.Refresh()
		ui.diskCheck.Refresh()
		ui.speedCheck.Refresh()
		ui.unlockCheck.Refresh()
		ui.routeCheck.Refresh()
	})

	buttonRow := container.NewHBox(selectAllBtn, deselectAllBtn)

	content := container.NewVBox(
		ui.basicCheck,
		ui.cpuCheck,
		ui.memoryCheck,
		ui.diskCheck,
		ui.speedCheck,
		ui.unlockCheck,
		ui.routeCheck,
		buttonRow,
	)

	return widget.NewCard("测试项目", "", content)
}

func (ui *TestUI) createConfigOptions() *widget.Card {
	// 语言选择
	ui.languageSelect = widget.NewSelect(
		[]string{"中文", "English"},
		func(value string) {},
	)
	ui.languageSelect.Selected = "中文"
	languageForm := container.NewBorder(
		nil, nil,
		widget.NewLabel("语言:"),
		nil,
		ui.languageSelect,
	)

	// CPU 测试方法
	ui.cpuMethodSelect = widget.NewSelect(
		[]string{"sysbench", "geekbench", "winsat"},
		func(value string) {},
	)
	ui.cpuMethodSelect.Selected = "sysbench"
	cpuMethodForm := container.NewBorder(
		nil, nil,
		widget.NewLabel("CPU方法:"),
		nil,
		ui.cpuMethodSelect,
	)

	// 线程模式
	ui.threadModeSelect = widget.NewSelect(
		[]string{"single", "multi"},
		func(value string) {},
	)
	ui.threadModeSelect.Selected = "multi"
	threadForm := container.NewBorder(
		nil, nil,
		widget.NewLabel("线程模式:"),
		nil,
		ui.threadModeSelect,
	)

	// 磁盘测试路径
	ui.diskPathEntry = widget.NewEntry()
	ui.diskPathEntry.SetPlaceHolder("/tmp 或留空自动检测")
	diskPathForm := container.NewBorder(
		nil, nil,
		widget.NewLabel("磁盘路径:"),
		nil,
		ui.diskPathEntry,
	)

	content := container.NewVBox(
		languageForm,
		cpuMethodForm,
		threadForm,
		diskPathForm,
	)

	return widget.NewCard("测试配置", "", content)
}

func (ui *TestUI) createControlButtons() *fyne.Container {
	ui.startButton = widget.NewButton("开始测试", ui.startTests)
	ui.startButton.Importance = widget.HighImportance

	ui.stopButton = widget.NewButton("停止", ui.stopTests)
	ui.stopButton.Disable()

	ui.clearButton = widget.NewButton("清空", ui.clearResults)

	// 使用VBox布局以适应小屏幕
	return container.NewVBox(
		ui.startButton,
		container.NewGridWithColumns(2,
			ui.stopButton,
			ui.clearButton,
		),
	)
}

func (ui *TestUI) createResultArea() *fyne.Container {
	// 状态标签
	ui.statusLabel = widget.NewLabel("就绪")

	// 进度条
	ui.progressBar = widget.NewProgressBar()
	ui.progressBar.Hide()

	// 结果文本框（多行，只读，可滚动）
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

func (ui *TestUI) startTests() {
	ui.mu.Lock()
	if ui.isRunning {
		ui.mu.Unlock()
		return
	}
	ui.isRunning = true
	ui.mu.Unlock()

	// 检查至少选择一项测试
	if !ui.hasSelectedTests() {
		dialog.ShowInformation("提示", "请至少选择一项测试！", ui.window)
		ui.mu.Lock()
		ui.isRunning = false
		ui.mu.Unlock()
		return
	}

	// 更新 UI 状态
	ui.startButton.Disable()
	ui.stopButton.Enable()
	ui.clearButton.Disable()
	ui.progressBar.Show()
	ui.statusLabel.SetText("测试运行中...")
	ui.resultText.SetText("")

	// 创建可取消的上下文
	ui.cancelCtx, ui.cancelFn = context.WithCancel(context.Background())

	// 在后台运行测试
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
		ui.speedCheck.Checked ||
		ui.unlockCheck.Checked ||
		ui.routeCheck.Checked
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

	ui.appendResult("========== GoECS Android 测试开始 ==========\n")
	ui.appendResult(fmt.Sprintf("测试时间: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	ui.appendResult(fmt.Sprintf("总测试项: %d\n\n", totalTests))

	// 基础信息测试
	if ui.basicCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "基础信息测试")
		ui.runBasicTest(language)
	}

	// CPU 测试
	if ui.cpuCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "CPU 性能测试")
		ui.runCPUTest()
	}

	// 内存测试
	if ui.memoryCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "内存性能测试")
		ui.runMemoryTest()
	}

	// 磁盘测试
	if ui.diskCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "磁盘性能测试")
		ui.runDiskTest()
	}

	// 网络测速
	if ui.speedCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "网络测速")
		ui.runSpeedTest()
	}

	// 流媒体解锁
	if ui.unlockCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "流媒体解锁测试")
		ui.runUnlockTest()
	}

	// 路由追踪
	if ui.routeCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "路由追踪测试")
		ui.runRouteTest()
	}

	if !ui.isCancelled() {
		ui.appendResult("\n========== 所有测试完成 ==========\n")
		ui.statusLabel.SetText("测试完成")
		dialog.ShowInformation("完成", "所有测试已完成！", ui.window)
	}
}

func (ui *TestUI) runBasicTest(language string) {
	ui.appendResult("\n--- 基础信息测试 ---\n")

	// 调用 utils 包获取基础信息
	ipv4, ipv6, result := utils.OnlyBasicsIpInfo(language)
	ui.appendResult(fmt.Sprintf("IPv4: %s\n", ipv4))
	ui.appendResult(fmt.Sprintf("IPv6: %s\n", ipv6))
	ui.appendResult(result + "\n")
}

func (ui *TestUI) runCPUTest() {
	ui.appendResult("\n--- CPU 性能测试 ---\n")

	method := ui.cpuMethodSelect.Selected
	threadMode := ui.threadModeSelect.Selected

	language := "zh"
	if ui.languageSelect.Selected == "English" {
		language = "en"
	}

	ui.appendResult(fmt.Sprintf("测试方法: %s\n", method))
	ui.appendResult(fmt.Sprintf("线程模式: %s\n", threadMode))

	// 调用 cputest 包
	realMethod, result := cputest.CpuTest(language, method, threadMode)
	ui.appendResult(fmt.Sprintf("实际使用方法: %s\n", realMethod))
	ui.appendResult(result + "\n")
}

func (ui *TestUI) runMemoryTest() {
	ui.appendResult("\n--- 内存性能测试 ---\n")

	language := "zh"
	if ui.languageSelect.Selected == "English" {
		language = "en"
	}

	// 调用 memorytest 包
	realMethod, result := memorytest.MemoryTest(language, "auto")
	ui.appendResult(fmt.Sprintf("测试方法: %s\n", realMethod))
	ui.appendResult(result + "\n")
}

func (ui *TestUI) runDiskTest() {
	ui.appendResult("\n--- 磁盘性能测试 ---\n")

	diskPath := ui.diskPathEntry.Text
	if diskPath == "" {
		diskPath = "/tmp"
	}

	language := "zh"
	if ui.languageSelect.Selected == "English" {
		language = "en"
	}

	ui.appendResult(fmt.Sprintf("测试路径: %s\n", diskPath))

	// 调用 disktest 包
	realMethod, result := disktest.DiskTest(language, "auto", diskPath, false, true)
	ui.appendResult(fmt.Sprintf("测试方法: %s\n", realMethod))
	ui.appendResult(result + "\n")
}

func (ui *TestUI) runSpeedTest() {
	ui.appendResult("\n--- 网络测速 ---\n")

	language := "zh"
	if ui.languageSelect.Selected == "English" {
		language = "en"
	}

	// 调用 speedtest 包
	speedtest.ShowHead(language)
	ui.appendResult("正在进行附近节点测速...\n")
	speedtest.NearbySP()
	ui.appendResult("测速完成\n")
}

func (ui *TestUI) runUnlockTest() {
	ui.appendResult("\n--- 流媒体解锁测试 ---\n")

	language := "zh"
	if ui.languageSelect.Selected == "English" {
		language = "en"
	}

	// 调用 unlocktest 包
	result := unlocktest.MediaTest(language)
	if result == "" {
		ui.appendResult("未检测到可用的网络连接\n")
	} else {
		ui.appendResult(result + "\n")
	}
}

func (ui *TestUI) runRouteTest() {
	ui.appendResult("\n--- 路由追踪测试 ---\n")
	ui.appendResult("路由追踪功能开发中...\n")
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
	if ui.speedCheck.Checked {
		count++
	}
	if ui.unlockCheck.Checked {
		count++
	}
	if ui.routeCheck.Checked {
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
	// 滚动到底部
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
