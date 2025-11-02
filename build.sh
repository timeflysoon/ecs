#!/bin/bash

set -e

BUILD_TYPE=${1:-"desktop"}

# 下载 ecs 二进制文件用于 embed
download_ecs_binary() {
    local target_os="$1"
    local target_arch="$2"
    
    echo "=========================================="
    echo "  下载 ECS 二进制文件用于 embed"
    echo "  目标平台: ${target_os}/${target_arch}"
    echo "=========================================="
    
    REPO="oneclickvirt/ecs"
    BINARIES_DIR="embedding/binaries"
    
    mkdir -p "$BINARIES_DIR"
    
    echo "获取最新版本信息..."
    LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest")
    ECS_VERSION=$(echo "$LATEST_RELEASE" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$ECS_VERSION" ]; then
        echo "警告: 无法获取最新版本，跳过下载"
        return
    fi
    
    echo "ECS 版本: $ECS_VERSION"
    
    local zipfile="goecs_${target_os}_${target_arch}.zip"
    local binname="goecs-${target_os}-${target_arch}"
    local exe_suffix=""
    
    # Windows 平台需要 .exe 后缀
    if [ "$target_os" = "windows" ]; then
        binname="${binname}.exe"
        exe_suffix=".exe"
    fi
    
    ZIP_URL="https://github.com/${REPO}/releases/download/${ECS_VERSION}/${zipfile}"
    BIN_PATH="${BINARIES_DIR}/${binname}"
    
    if [ -f "$BIN_PATH" ]; then
        echo "✓ ${binname} 已存在"
        return
    fi
    
    echo "下载 ${target_os}/${target_arch}..."
    
    TMP_ZIP="/tmp/${zipfile}"
    if curl -L -f -o "$TMP_ZIP" "$ZIP_URL"; then
        unzip -q -o "$TMP_ZIP" -d /tmp/
        
        local extracted_file="/tmp/goecs${exe_suffix}"
        if [ -f "$extracted_file" ]; then
            mv "$extracted_file" "$BIN_PATH"
            chmod +x "$BIN_PATH"
            echo "✓ 下载并解压成功: ${binname}"
            echo "  文件大小: $(du -h "$BIN_PATH" | cut -f1)"
        else
            echo "✗ 解压失败: 未找到可执行文件 goecs${exe_suffix}"
        fi
        
        rm -f "$TMP_ZIP"
    else
        echo "✗ 下载失败: ${zipfile}"
    fi
    
    echo ""
}

# 检查 Fyne CLI 是否安装
check_fyne_cli() {
    if ! command -v fyne &> /dev/null; then
        echo "正在安装 Fyne CLI..."
        go install fyne.io/fyne/v2/cmd/fyne@latest
        if [ $? -ne 0 ]; then
            echo "Fyne CLI 安装失败"
            exit 1
        fi
        echo "Fyne CLI 安装成功"
    else
        echo "Fyne CLI 已安装"
    fi
}

# 桌面端构建（用于快速测试）
build_desktop() {
    # 检测当前平台
    local current_os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local current_arch=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
    
    echo "检测到当前平台: ${current_os}/${current_arch}"
    
    download_ecs_binary "$current_os" "$current_arch"
    
    go build -ldflags="-checklinkname=0 -s -w" -o goecs-desktop .
    
    if [ $? -eq 0 ]; then
        echo "桌面端编译成功！"
        ls -lh goecs-desktop
    else
        echo "桌面端编译失败"
        exit 1
    fi
}

# 获取版本信息
get_version() {
    VERSION="v0.0.1-$(date +%Y%m%d)-$(git rev-parse --short HEAD 2>/dev/null || echo 'dev')"
    echo "$VERSION"
}

# macOS 构建
build_macos() {
    VERSION=$(get_version)
    echo "构建 macOS 版本 - 版本: $VERSION"
    
    mkdir -p .build
    
    echo ""
    echo "构建 macOS ARM64 版本..."
    download_ecs_binary "darwin" "arm64"
    fyne package -os darwin/arm64 -name goecs
    if [ -f goecs.app ] || [ -d goecs.app ]; then
        tar -czf .build/goecs-macos-arm64-${VERSION}.tar.gz goecs.app
        rm -rf goecs.app
        echo "macOS ARM64 构建成功"
    else
        echo "macOS ARM64 构建失败"
        exit 1
    fi
    
    echo ""
    echo "构建 macOS AMD64 版本..."
    download_ecs_binary "darwin" "amd64"
    fyne package -os darwin/amd64 -name goecs
    if [ -f goecs.app ] || [ -d goecs.app ]; then
        tar -czf .build/goecs-macos-amd64-${VERSION}.tar.gz goecs.app
        rm -rf goecs.app
        echo "macOS AMD64 构建成功"
    else
        echo "macOS AMD64 构建失败"
        exit 1
    fi
}

# Windows 构建
build_windows() {
    VERSION=$(get_version)
    echo "构建 Windows 版本 - 版本: $VERSION"
    
    mkdir -p .build
    
    echo ""
    echo "构建 Windows ARM64 版本..."
    download_ecs_binary "windows" "arm64"
    fyne package -os windows/arm64 -name goecs
    if [ -f goecs.exe ]; then
        mv goecs.exe .build/goecs-windows-arm64-${VERSION}.exe
        echo "Windows ARM64 构建成功"
    else
        echo "Windows ARM64 构建失败"
        exit 1
    fi
    
    echo ""
    echo "构建 Windows AMD64 版本..."
    download_ecs_binary "windows" "amd64"
    fyne package -os windows/amd64 -name goecs
    if [ -f goecs.exe ]; then
        mv goecs.exe .build/goecs-windows-amd64-${VERSION}.exe
        echo "Windows AMD64 构建成功"
    else
        echo "Windows AMD64 构建失败"
        exit 1
    fi
}

# Linux 构建
build_linux() {
    VERSION=$(get_version)
    echo "构建 Linux 版本 - 版本: $VERSION"
    
    mkdir -p .build
    
    echo ""
    echo "构建 Linux ARM64 版本..."
    download_ecs_binary "linux" "arm64"
    fyne package -os linux/arm64 -name goecs
    if [ -f goecs.tar.xz ]; then
        mv goecs.tar.xz .build/goecs-linux-arm64-${VERSION}.tar.xz
        echo "Linux ARM64 构建成功"
    else
        echo "Linux ARM64 构建失败"
        exit 1
    fi
    
    echo ""
    echo "构建 Linux AMD64 版本..."
    download_ecs_binary "linux" "amd64"
    fyne package -os linux/amd64 -name goecs
    if [ -f goecs.tar.xz ]; then
        mv goecs.tar.xz .build/goecs-linux-amd64-${VERSION}.tar.xz
        echo "Linux AMD64 构建成功"
    else
        echo "Linux AMD64 构建失败"
        exit 1
    fi
}

# 准备 Android 的 jniLibs
prepare_android_jnilibs() {
    echo "=========================================="
    echo "  准备 Android JNI 库"
    echo "=========================================="
    
    # 下载 ARM64 版本的 ECS 二进制
    echo "下载 ARM64 ECS 二进制..."
    download_ecs_binary "linux" "arm64"
    
    # 下载 x86_64 版本的 ECS 二进制
    echo "下载 x86_64 ECS 二进制..."
    download_ecs_binary "linux" "amd64"
    
    # 创建 jniLibs 目录结构
    mkdir -p jniLibs/arm64-v8a
    mkdir -p jniLibs/x86_64
    
    # 复制二进制文件并重命名为 .so 库格式
    if [ -f "embedding/binaries/goecs-linux-arm64" ]; then
        cp "embedding/binaries/goecs-linux-arm64" "jniLibs/arm64-v8a/libgoecs.so"
        chmod 755 "jniLibs/arm64-v8a/libgoecs.so"
        echo "✓ ARM64 库已准备: jniLibs/arm64-v8a/libgoecs.so"
        echo "  文件大小: $(du -h jniLibs/arm64-v8a/libgoecs.so | cut -f1)"
    else
        echo "✗ 错误: 未找到 ARM64 ECS 二进制文件"
        exit 1
    fi
    
    if [ -f "embedding/binaries/goecs-linux-amd64" ]; then
        cp "embedding/binaries/goecs-linux-amd64" "jniLibs/x86_64/libgoecs.so"
        chmod 755 "jniLibs/x86_64/libgoecs.so"
        echo "✓ x86_64 库已准备: jniLibs/x86_64/libgoecs.so"
        echo "  文件大小: $(du -h jniLibs/x86_64/libgoecs.so | cut -f1)"
    else
        echo "✗ 错误: 未找到 x86_64 ECS 二进制文件"
        exit 1
    fi
    
    echo ""
    echo "JNI 库准备完成！"
    echo ""
}

# Android 构建
build_android() {
    VERSION=$(get_version)
    echo "构建 Android 版本 - 版本: $VERSION"
    
    if [ -z "$ANDROID_NDK_HOME" ]; then
        echo "请设置 Android NDK 路径，例如："
        echo "export ANDROID_NDK_HOME=/path/to/android-ndk"
        exit 1
    fi
    
    echo "Android NDK: $ANDROID_NDK_HOME"
    
    mkdir -p .build
    
    # 准备 JNI 库
    prepare_android_jnilibs
    
    echo ""
    echo "构建 Android APK..."
    
    # 构建包含所有架构的 APK
    fyne package -os android -appID com.oneclickvirt.goecs -appVersion "$VERSION"
    
    if [ -f *.apk ]; then
        mv *.apk .build/goecs-android-${VERSION}.apk
        echo "Android APK 构建成功"
    else
        echo "Android APK 构建失败"
        exit 1
    fi
    
    echo ""
    echo "=========================================="
    echo "  Android 构建完成"
    echo "=========================================="
    ls -lh .build/*.apk
}

# 主流程
case "$BUILD_TYPE" in
    "desktop")
        build_desktop
        ;;
    "android")
        check_fyne_cli
        build_android
        ;;
    "macos")
        check_fyne_cli
        build_macos
        ;;
    "windows")
        check_fyne_cli
        build_windows
        ;;
    "linux")
        check_fyne_cli
        build_linux
        ;;
    "all")
        build_desktop
        echo ""
        check_fyne_cli
        echo ""
        build_macos
        echo ""
        build_windows
        echo ""
        build_linux
        echo ""
        build_android
        ;;
    *)
        echo "用法: $0 [desktop|android|macos|windows|linux|all]"
        echo ""
        echo "  desktop - 构建桌面端应用（默认，用于快速测试）"
        echo "  android - 构建 Android APK (arm64 + x86_64)"
        echo "  macos   - 构建 macOS 应用 (arm64 + amd64)"
        echo "  windows - 构建 Windows 应用 (arm64 + amd64)"
        echo "  linux   - 构建 Linux 应用 (arm64 + amd64)"
        echo "  all     - 构建所有平台"
        exit 1
        ;;
esac

echo ""
echo "=========================================="
echo "  所有构建任务完成"
echo "=========================================="
if [ -d .build ]; then
    echo ""
    echo "构建产物:"
    ls -lh .build/
    echo ""
    echo "总大小:"
    du -sh .build/
fi
