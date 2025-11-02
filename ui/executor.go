package ui

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/oneclickvirt/basics/utils"
	"github.com/oneclickvirt/ecs/cputest"
	"github.com/oneclickvirt/ecs/disktest"
	"github.com/oneclickvirt/ecs/memorytest"
	"github.com/oneclickvirt/ecs/unlocktest"
	ecsutils "github.com/oneclickvirt/ecs/utils"
	"github.com/oneclickvirt/pingtest/pt"
	"github.com/oneclickvirt/portchecker/email"
)

const ecsVersion = "v0.1.93"

type CommandExecutor struct {
	outputCallback func(string)
}

func NewCommandExecutor(outputCallback func(string)) *CommandExecutor {
	return &CommandExecutor{
		outputCallback: outputCallback,
	}
}

func (e *CommandExecutor) Execute(selectedOptions map[string]bool, language string) error {
	// 设置测试选项（模仿 main 函数中的变量）
	basicStatus := selectedOptions["basic"]
	cpuTestStatus := selectedOptions["cpu"]
	memoryTestStatus := selectedOptions["memory"]
	diskTestStatus := selectedOptions["disk"]
	utTestStatus := selectedOptions["unlock"]
	securityTestStatus := selectedOptions["security"]
	emailTestStatus := selectedOptions["email"]
	speedTestStatus := selectedOptions["speed"]

	// 检查网络连接
	preCheck := utils.CheckPublicAccess(3 * time.Second)

	// 初始化变量（完全模仿 main 函数）
	var (
		wg1, wg2                                      sync.WaitGroup
		basicInfo, securityInfo, emailInfo, mediaInfo string
		output, tempOutput                            string
		outputMutex                                   sync.Mutex
	)
	startTime := time.Now()
	uploadDone := make(chan bool, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 重定向输出到回调
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan bool)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if n > 0 && e.outputCallback != nil {
				e.outputCallback(string(buf[:n]))
			}
			if err != nil {
				if err == io.EOF {
					break
				}
			}
		}
		done <- true
	}()

	// 执行测试（模仿 runChineseTests/runEnglishTests 的逻辑）
	// 基础信息测试
	if basicStatus || securityTestStatus {
		outputMutex.Lock()
		ecsutils.PrintHead(language, 82, ecsVersion)
		if basicStatus {
			if language == "zh" {
				ecsutils.PrintCenteredTitle("系统基础信息", 82)
			} else {
				ecsutils.PrintCenteredTitle("System-Basic-Information", 82)
			}
		}
		if preCheck.Connected {
			_, _, basicInfo, securityInfo, _ = ecsutils.BasicsAndSecurityCheck(language, "ipv4", securityTestStatus)
		}
		if basicStatus {
			fmt.Printf("%s", basicInfo)
		}
		outputMutex.Unlock()
	}

	// CPU测试
	if cpuTestStatus {
		outputMutex.Lock()
		realTestMethod, res := cputest.CpuTest(language, "sysbench", "multi")
		if language == "zh" {
			ecsutils.PrintCenteredTitle(fmt.Sprintf("CPU测试-通过%s测试", realTestMethod), 82)
		} else {
			ecsutils.PrintCenteredTitle(fmt.Sprintf("CPU-Test--%s-Method", realTestMethod), 82)
		}
		fmt.Print(res)
		outputMutex.Unlock()
	}

	// 内存测试
	if memoryTestStatus {
		outputMutex.Lock()
		realTestMethod, res := memorytest.MemoryTest(language, "auto")
		if language == "zh" {
			ecsutils.PrintCenteredTitle(fmt.Sprintf("内存测试-通过%s测试", realTestMethod), 82)
		} else {
			ecsutils.PrintCenteredTitle(fmt.Sprintf("Memory-Test--%s-Method", realTestMethod), 82)
		}
		fmt.Print(res)
		outputMutex.Unlock()
	}

	// 磁盘测试
	if diskTestStatus {
		outputMutex.Lock()
		realTestMethod, res := disktest.DiskTest(language, "fio", "", false, true)
		if language == "zh" {
			ecsutils.PrintCenteredTitle(fmt.Sprintf("硬盘测试-通过%s测试", realTestMethod), 82)
		} else {
			ecsutils.PrintCenteredTitle(fmt.Sprintf("Disk-Test--%s-Method", realTestMethod), 82)
		}
		fmt.Print(res)
		outputMutex.Unlock()
	}

	// 流媒体解锁测试
	if utTestStatus && preCheck.Connected {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			mediaInfo = unlocktest.MediaTest(language)
		}()
	}

	// 邮件端口测试
	if emailTestStatus && preCheck.Connected {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			emailInfo = email.EmailCheck()
		}()
	}

	// 显示流媒体解锁结果
	if utTestStatus && preCheck.Connected {
		wg1.Wait()
		outputMutex.Lock()
		if language == "zh" {
			ecsutils.PrintCenteredTitle("跨国流媒体解锁", 82)
		} else {
			ecsutils.PrintCenteredTitle("Cross-Border-Streaming-Media-Unlock", 82)
		}
		fmt.Printf("%s", mediaInfo)
		outputMutex.Unlock()
	}

	// 显示IP质量检测结果
	if securityTestStatus && preCheck.Connected {
		outputMutex.Lock()
		if language == "zh" {
			ecsutils.PrintCenteredTitle("IP质量检测", 82)
		} else {
			ecsutils.PrintCenteredTitle("IP-Quality-Check", 82)
		}
		fmt.Printf("%s", securityInfo)
		outputMutex.Unlock()
	}

	// 显示邮件端口测试结果
	if emailTestStatus && preCheck.Connected {
		wg2.Wait()
		outputMutex.Lock()
		if language == "zh" {
			ecsutils.PrintCenteredTitle("邮件端口检测", 82)
		} else {
			ecsutils.PrintCenteredTitle("Email-Port-Check", 82)
		}
		fmt.Println(emailInfo)
		outputMutex.Unlock()
	}

	// 速度测试（简化版）
	if speedTestStatus && preCheck.Connected {
		// 这里可以添加速度测试，但需要导入 speedtest 包
		_ = pt.PingTest // 避免未使用的导入
	}

	// 打印时间信息
	outputMutex.Lock()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	currentTime := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
	ecsutils.PrintCenteredTitle("", 82)
	if language == "zh" {
		fmt.Printf("花费          : %d 分 %d 秒\n", minutes, seconds)
		fmt.Printf("时间          : %s\n", currentTime)
	} else {
		fmt.Printf("Cost    Time          : %d min %d sec\n", minutes, seconds)
		fmt.Printf("Current Time          : %s\n", currentTime)
	}
	ecsutils.PrintCenteredTitle("", 82)
	outputMutex.Unlock()

	// 清理
	_ = uploadDone
	_ = sig
	_ = output
	_ = tempOutput

	// 恢复输出
	w.Close()
	<-done
	os.Stdout = oldStdout

	return nil
}
