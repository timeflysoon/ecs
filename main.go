package main

import (
	"flag"
	"fmt"
	"os"

	"fyne.io/fyne/v2/app"
)

var (
	// UI 模式标志
	guiMode bool

	// 测试选项标志
	basicTest  bool
	cpuTest    bool
	memoryTest bool
	diskTest   bool
	speedTest  bool
	unlockTest bool
	routeTest  bool
	allTests   bool

	// 配置选项
	language   string
	cpuMethod  string
	threadMode string
	diskPath   string
	diskMethod string

	// 其他选项
	showVersion bool
	showHelp    bool
)

func init() {
	// UI 模式
	flag.BoolVar(&guiMode, "gui", false, "启动图形界面模式 / Launch GUI mode")

	// 测试选项
	flag.BoolVar(&basicTest, "basic", false, "基础信息测试 / Basic info test")
	flag.BoolVar(&cpuTest, "cpu", false, "CPU 性能测试 / CPU performance test")
	flag.BoolVar(&memoryTest, "memory", false, "内存性能测试 / Memory performance test")
	flag.BoolVar(&diskTest, "disk", false, "磁盘性能测试 / Disk performance test")
	flag.BoolVar(&speedTest, "speed", false, "网络测速 / Network speed test")
	flag.BoolVar(&unlockTest, "unlock", false, "流媒体解锁测试 / Media unlock test")
	flag.BoolVar(&routeTest, "route", false, "路由追踪测试 / Route trace test")
	flag.BoolVar(&allTests, "all", false, "运行所有测试 / Run all tests")

	// 配置选项
	flag.StringVar(&language, "lang", "zh", "语言: zh/en / Language: zh/en")
	flag.StringVar(&cpuMethod, "cpu-method", "sysbench", "CPU测试方法: sysbench/geekbench/winsat")
	flag.StringVar(&threadMode, "thread", "multi", "线程模式: single/multi")
	flag.StringVar(&diskPath, "disk-path", "", "磁盘测试路径 / Disk test path")
	flag.StringVar(&diskMethod, "disk-method", "auto", "磁盘测试方法: fio/dd/auto")

	// 其他选项
	flag.BoolVar(&showVersion, "version", false, "显示版本信息 / Show version")
	flag.BoolVar(&showVersion, "v", false, "显示版本信息 / Show version")
	flag.BoolVar(&showHelp, "help", false, "显示帮助信息 / Show help")
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息 / Show help")
}

func main() {
	flag.Parse()

	// 显示版本信息
	if showVersion {
		fmt.Println("GoECS Android v1.0.0")
		fmt.Println("Based on github.com/oneclickvirt/ecs v0.1.91")
		os.Exit(0)
	}

	// 显示帮助信息
	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// 如果指定了 GUI 模式或没有任何参数，启动 UI
	if guiMode || (!hasAnyTest() && flag.NFlag() == 0) {
		runGUIMode()
	} else {
		runCLIMode()
	}
}

func runGUIMode() {
	myApp := app.NewWithID("com.oneclickvirt.goecs")
	myApp.Settings().SetTheme(&customTheme{})

	ui := NewTestUI(myApp)
	ui.window.ShowAndRun()
}

func runCLIMode() {
	fmt.Println("========== GoECS Android CLI 模式 ==========")

	// 如果指定了 -all，启用所有测试
	if allTests {
		basicTest = true
		cpuTest = true
		memoryTest = true
		diskTest = true
		speedTest = true
		unlockTest = true
		routeTest = true
	}

	// 执行测试
	runner := NewCLIRunner(language, cpuMethod, threadMode, diskPath, diskMethod)
	runner.Run(basicTest, cpuTest, memoryTest, diskTest, speedTest, unlockTest, routeTest)
}

func hasAnyTest() bool {
	return basicTest || cpuTest || memoryTest || diskTest || speedTest || unlockTest || routeTest || allTests
}

func printHelp() {
	fmt.Println(`
GoECS Android - 服务器性能测试工具 (Android 版本)

用法 / Usage:
  goecs-android [选项] [测试项]

模式 / Modes:
  -gui              启动图形界面模式（默认）
                    Launch GUI mode (default)

测试项 / Tests:
  -basic            基础信息测试
  -cpu              CPU 性能测试
  -memory           内存性能测试
  -disk             磁盘性能测试
  -speed            网络测速
  -unlock           流媒体解锁测试
  -route            路由追踪测试
  -all              运行所有测试

配置选项 / Configuration:
  -lang string      语言: zh/en (默认: zh)
  -cpu-method       CPU测试方法: sysbench/geekbench/winsat (默认: sysbench)
  -thread           线程模式: single/multi (默认: multi)
  -disk-path        磁盘测试路径 (默认: 自动检测)
  -disk-method      磁盘测试方法: fio/dd/auto (默认: auto)

其他选项 / Other:
  -version, -v      显示版本信息
  -help, -h         显示此帮助信息

示例 / Examples:
  # 启动 GUI 模式
  goecs-android
  goecs-android -gui

  # 运行所有测试（CLI 模式）
  goecs-android -all

  # 运行指定测试
  goecs-android -basic -cpu -memory
  
  # 指定配置运行测试
  goecs-android -cpu -cpu-method sysbench -thread multi
  goecs-android -disk -disk-path /tmp -disk-method fio
  
  # 英文环境运行
  goecs-android -all -lang en

更多信息 / More info:
  GitHub: https://github.com/oneclickvirt/ecs
  分支 / Branch: android-app
`)
}
