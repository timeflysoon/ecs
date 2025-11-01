package ui

import (
	"regexp"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TerminalOutput 是一个类似终端的输出组件
type TerminalOutput struct {
	widget.Entry
	mu       sync.Mutex
	lines    []string
	maxLines int
}

// NewTerminalOutput 创建新的终端输出组件
func NewTerminalOutput() *TerminalOutput {
	terminal := &TerminalOutput{
		lines:    make([]string, 0),
		maxLines: 10000, // 最大行数限制
	}
	terminal.ExtendBaseWidget(terminal)
	terminal.MultiLine = true
	terminal.Wrapping = fyne.TextWrapWord
	terminal.Disable() // 禁用编辑
	return terminal
}

// AppendText 追加文本到终端
func (t *TerminalOutput) AppendText(text string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 移除ANSI颜色代码
	cleanText := t.stripANSI(text)

	// 分割成行
	newLines := strings.Split(cleanText, "\n")

	// 如果最后一行是空的，移除它（避免额外的空行）
	if len(newLines) > 0 && newLines[len(newLines)-1] == "" {
		newLines = newLines[:len(newLines)-1]
	}

	t.lines = append(t.lines, newLines...)

	// 限制最大行数
	if len(t.lines) > t.maxLines {
		t.lines = t.lines[len(t.lines)-t.maxLines:]
	}

	// 更新显示
	t.SetText(strings.Join(t.lines, "\n"))

	// 自动滚动到底部
	t.CursorRow = len(t.lines)
}

// Clear 清空终端内容
func (t *TerminalOutput) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.lines = make([]string, 0)
	t.SetText("")
}

// SetText 设置完整文本（覆盖现有内容）
func (t *TerminalOutput) SetFullText(text string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	cleanText := t.stripANSI(text)
	t.lines = strings.Split(cleanText, "\n")

	// 限制最大行数
	if len(t.lines) > t.maxLines {
		t.lines = t.lines[len(t.lines)-t.maxLines:]
	}

	t.SetText(strings.Join(t.lines, "\n"))
	t.CursorRow = len(t.lines)
}

// stripANSI 移除ANSI转义序列
func (t *TerminalOutput) stripANSI(text string) string {
	ansiRegex := regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(text, "")
}

// GetText 获取当前文本内容
func (t *TerminalOutput) GetText() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return strings.Join(t.lines, "\n")
}
