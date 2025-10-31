# GoECS - 跨平台测试工具

[![Build All UI APP](https://github.com/oneclickvirt/ecs/actions/workflows/build-all.yml/badge.svg)](https://github.com/oneclickvirt/ecs/actions/workflows/build-all.yml)

一个基于 Fyne 框架的跨平台测试工具，支持 Android、macOS、Windows 和 Linux。

## 支持的平台

### Android
- Android 7.0 (API Level 24) 或更高版本
- 建议 Android 13 (API Level 33) 以获得最佳体验
- 支持架构：ARM64、x86_64

### macOS
- macOS 11.0 或更高版本
- 支持架构：Apple Silicon (ARM64)、Intel (AMD64)

### Windows
- Windows 10 或更高版本
- 支持架构：ARM64、AMD64

### Linux
- 主流 Linux 发行版
- 支持架构：ARM64、AMD64

## 本地构建

### 前置要求

1. Go 1.25.3
2. Android SDK
3. Android NDK 25.2.9519653
4. JDK 17+

### 环境配置

```bash
# 设置 Android NDK 路径
export ANDROID_NDK_HOME=/path/to/android-ndk

# 安装 Fyne CLI
go install fyne.io/fyne/v2/cmd/fyne@latest
```

### 构建命令

```bash
# 构建桌面端（用于快速测试）
./build.sh desktop

# 构建 Android APK (arm64 + x86_64)
./build.sh android

# 构建 macOS 应用 (arm64 + amd64)
./build.sh macos

# 构建 Windows 应用 (arm64 + amd64)
./build.sh windows

# 构建 Linux 应用 (arm64 + amd64)
./build.sh linux

# 构建所有平台
./build.sh all
```

构建产物将输出到 `.build/` 目录。

### 构建产物说明

- **Android**: `.apk` 文件
  - `goecs-android-arm64-*.apk` - ARM64 版本（真机）
  - `goecs-android-x86_64-*.apk` - x86_64 版本（模拟器）

- **macOS**: `.tar.gz` 压缩包（包含 `.app` 应用）
  - `goecs-macos-arm64-*.tar.gz` - Apple Silicon 版本
  - `goecs-macos-amd64-*.tar.gz` - Intel 版本

- **Windows**: `.exe` 可执行文件
  - `goecs-windows-arm64-*.exe` - ARM64 版本
  - `goecs-windows-amd64-*.exe` - AMD64 版本

- **Linux**: `.tar.gz` 压缩包（包含可执行文件）
  - `goecs-linux-arm64-*.tar.gz` - ARM64 版本
  - `goecs-linux-amd64-*.tar.gz` - AMD64 版本

## 开发

```bash
# 克隆仓库
git clone https://github.com/oneclickvirt/ecs.git
cd ecs

# 切换到 Android 开发分支
git checkout android-app

# 安装依赖
go mod download

# 运行桌面版本（用于开发测试）
go run -ldflags="-checklinkname=0" .
```