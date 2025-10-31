package main

import (
	"fmt"
	"time"

	"github.com/oneclickvirt/ecs/cputest"
	"github.com/oneclickvirt/ecs/disktest"
	"github.com/oneclickvirt/ecs/memorytest"
	"github.com/oneclickvirt/ecs/speedtest"
	"github.com/oneclickvirt/ecs/unlocktest"
	"github.com/oneclickvirt/ecs/utils"
)

// CLIRunner CLI 模式的测试运行器
type CLIRunner struct {
	language   string
	cpuMethod  string
	threadMode string
	diskPath   string
	diskMethod string
}

// NewCLIRunner 创建新的 CLI 运行器
func NewCLIRunner(language, cpuMethod, threadMode, diskPath, diskMethod string) *CLIRunner {
	if diskPath == "" {
		diskPath = "/tmp"
	}
	return &CLIRunner{
		language:   language,
		cpuMethod:  cpuMethod,
		threadMode: threadMode,
		diskPath:   diskPath,
		diskMethod: diskMethod,
	}
}

// Run 运行测试
func (r *CLIRunner) Run(basic, cpu, memory, disk, speed, unlock, route bool) {
	startTime := time.Now()

	fmt.Printf("\n测试开始时间: %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Println("=" + repeatString("=", 80))

	testCount := 0

	// 基础信息测试
	if basic {
		testCount++
		r.runBasicTest()
	}

	// CPU 测试
	if cpu {
		testCount++
		r.runCPUTest()
	}

	// 内存测试
	if memory {
		testCount++
		r.runMemoryTest()
	}

	// 磁盘测试
	if disk {
		testCount++
		r.runDiskTest()
	}

	// 网络测速
	if speed {
		testCount++
		r.runSpeedTest()
	}

	// 流媒体解锁
	if unlock {
		testCount++
		r.runUnlockTest()
	}

	// 路由追踪
	if route {
		testCount++
		r.runRouteTest()
	}

	// 总结
	elapsed := time.Since(startTime)
	fmt.Println("\n" + repeatString("=", 80))
	fmt.Printf("测试完成! 共执行 %d 项测试，耗时: %v\n", testCount, elapsed)
	fmt.Println(repeatString("=", 80))
}

func (r *CLIRunner) runBasicTest() {
	fmt.Println("\n[1/N] 基础信息测试")
	fmt.Println(repeatString("-", 80))

	ipv4, ipv6, result := utils.OnlyBasicsIpInfo(r.language)
	fmt.Printf("IPv4: %s\n", ipv4)
	fmt.Printf("IPv6: %s\n", ipv6)
	fmt.Println(result)
}

func (r *CLIRunner) runCPUTest() {
	fmt.Println("\n[N/N] CPU 性能测试")
	fmt.Println(repeatString("-", 80))
	fmt.Printf("测试方法: %s\n", r.cpuMethod)
	fmt.Printf("线程模式: %s\n", r.threadMode)

	realMethod, result := cputest.CpuTest(r.language, r.cpuMethod, r.threadMode)
	fmt.Printf("实际使用方法: %s\n", realMethod)
	fmt.Println(result)
}

func (r *CLIRunner) runMemoryTest() {
	fmt.Println("\n[N/N] 内存性能测试")
	fmt.Println(repeatString("-", 80))

	realMethod, result := memorytest.MemoryTest(r.language, "auto")
	fmt.Printf("测试方法: %s\n", realMethod)
	fmt.Println(result)
}

func (r *CLIRunner) runDiskTest() {
	fmt.Println("\n[N/N] 磁盘性能测试")
	fmt.Println(repeatString("-", 80))
	fmt.Printf("测试路径: %s\n", r.diskPath)
	fmt.Printf("测试方法: %s\n", r.diskMethod)

	realMethod, result := disktest.DiskTest(r.language, r.diskMethod, r.diskPath, false, true)
	fmt.Printf("实际使用方法: %s\n", realMethod)
	fmt.Println(result)
}

func (r *CLIRunner) runSpeedTest() {
	fmt.Println("\n[N/N] 网络测速")
	fmt.Println(repeatString("-", 80))

	speedtest.ShowHead(r.language)
	fmt.Println("正在进行附近节点测速...")
	speedtest.NearbySP()
	fmt.Println("测速完成")
}

func (r *CLIRunner) runUnlockTest() {
	fmt.Println("\n[N/N] 流媒体解锁测试")
	fmt.Println(repeatString("-", 80))

	result := unlocktest.MediaTest(r.language)
	if result == "" {
		fmt.Println("未检测到可用的网络连接")
	} else {
		fmt.Println(result)
	}
}

func (r *CLIRunner) runRouteTest() {
	fmt.Println("\n[N/N] 路由追踪测试")
	fmt.Println(repeatString("-", 80))
	fmt.Println("路由追踪功能开发中...")
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
