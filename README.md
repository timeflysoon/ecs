# GoECS Android App

GoECS 服务器性能测试工具的 Android 版本。支持 **图形界面模式** 和 **命令行模式**。

## ✨ 特性

- ✅ **双模式运行**: 图形界面 (GUI) / 命令行 (CLI)
- ✅ **7 种测试项目**
  - 基础信息测试 (系统信息、IP 信息)
  - CPU 性能测试 (sysbench/geekbench/winsat)
  - 内存性能测试
  - 磁盘性能测试 (fio/dd)
  - 网络测速
  - 流媒体解锁测试
  - 路由追踪测试
- ✅ **后台执行 + 实时进度**
- ✅ **结果导出**
- ✅ **完整的参数支持**

## 快速开始

### 运行应用

```bash
# 直接运行
./goecs-android

# 显示帮助
./goecs-android --help

# 显示版本
./goecs-android -version
```

## 📖 命令行参数

### 模式选择
- `-gui`: 启动图形界面模式（默认）

### 测试项
- `-basic`: 基础信息测试
- `-cpu`: CPU 性能测试
- `-memory`: 内存性能测试
- `-disk`: 磁盘性能测试
- `-speed`: 网络测速
- `-unlock`: 流媒体解锁测试
- `-route`: 路由追踪测试
- `-all`: 运行所有测试

### 配置选项
- `-lang string`: 语言 (zh/en，默认: zh)
- `-cpu-method string`: CPU 测试方法 (sysbench/geekbench/winsat，默认: sysbench)
- `-thread string`: 线程模式 (single/multi，默认: multi)
- `-disk-path string`: 磁盘测试路径（默认: 自动检测）
- `-disk-method string`: 磁盘测试方法 (fio/dd/auto，默认: auto)

### 其他
- `-version, -v`: 显示版本信息
- `-help, -h`: 显示帮助信息

## 🔨 本地开发

### 桌面调试（快速测试 UI）
```bash
go run .
```

## 本地开发

### 前置要求

- Go 1.21+
- Fyne v2.4.5+
- 用于 Android 构建：Android SDK + NDK

### macOS 上测试

```bash
# 安装依赖
go mod download

# 运行桌面版本（用于开发测试）
go run .

# 或编译后运行
go build -o goecs-android .
./goecs-android
```

### 构建 Android APK

#### 方法 1: 使用 Fyne CLI（本地构建）

```bash
# 安装 Fyne CLI
go install fyne.io/fyne/v2/cmd/fyne@latest

# 构建 APK（多架构）
mkdir -p .build

# ARM64 架构（主流 Android 设备）
fyne package -os android -appID com.oneclickvirt.goecs -name GoECS
mv GoECS.apk .build/goecs-android-arm64.apk

# x86_64 架构（模拟器）
ANDROID_ARCH=x86_64 fyne package -os android -appID com.oneclickvirt.goecs -name GoECS
mv GoECS.apk .build/goecs-android-x86_64.apk
```

#### 方法 2: 使用 GitHub Actions（推荐）

手动触发构建：

1. 访问 GitHub 仓库的 Actions 页面
2. 选择 "Build Android APK" workflow
3. 点击 "Run workflow" 按钮
4. 等待构建完成

构建成功后：
- APK 文件会自动提交到 `android-app` 分支的 `.build/` 目录
- 同时在 Actions 的 Artifacts 中也可以下载

文件命名格式：
- `goecs-android-arm64-v1.0.0-YYYYMMDD-{hash}.apk` - 真机使用
- `goecs-android-x86_64-v1.0.0-YYYYMMDD-{hash}.apk` - 模拟器使用

## 技术栈

- **UI 框架**: [Fyne](https://fyne.io/) v2.4.5
- **核心库**: [github.com/oneclickvirt/ecs](https://github.com/oneclickvirt/ecs) v0.1.91
- **语言**: Go 1.21+
- **目标平台**: Android 7.0+ (API Level 24+)
- **支持架构**: ARM64, x86_64

## 常见问题

### Q: 为什么界面显示方块字符？

A: 需要确保 Android 系统有中文字体支持。本应用使用 Android 系统默认字体，应该能正常显示中文。如果问题依然存在，请检查设备的语言设置。

### Q: 如何在 macOS 上测试？

A: 使用 `go run .` 或 `go build` 编译后直接运行。macOS 版本仅用于开发测试，最终产品为 Android APK。

### Q: APK 文件在哪里？

A: 本地构建后在 `.build/` 目录；GitHub Actions 构建后会自动提交到 android-app 分支的 `.build/` 目录，也可从 Artifacts 下载。

### Q: 支持哪些 Android 设备？

A: Android 7.0 (API 24) 及以上版本。ARM64 APK 适用于大多数现代设备，x86_64 APK 适用于模拟器。

### Q: 如何手动触发 APK 构建？

A: 访问 GitHub 仓库的 Actions 页面，选择 "Build Android APK" workflow，点击 "Run workflow" 按钮即可。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

遵循主项目 https://github.com/oneclickvirt/ecs 的许可证。

## 相关链接

- 主项目: https://github.com/oneclickvirt/ecs
- Android 分支: https://github.com/oneclickvirt/ecs/tree/android-app

## 分支说明

这是一个**孤儿分支**（orphan branch），与 master 分支完全独立：
- 没有 master 的提交历史
- 根目录直接是应用代码
- 使用远程依赖而非本地引用
- 独立的 CI/CD 流程
