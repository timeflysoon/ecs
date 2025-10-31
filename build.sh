#!/bin/bash

set -e

BUILD_TYPE=${1:-"desktop"}

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
    go build -ldflags="-checklinkname=0 -s -w" -o goecs-desktop .
    
    if [ $? -eq 0 ]; then
        echo "桌面端编译成功！"
        ls -lh goecs-desktop
    else
        echo "桌面端编译失败"
        exit 1
    fi
}

# Android 构建
build_android() {
    echo ""
    
    # 检查必要的环境变量
    if [ -z "$ANDROID_NDK_HOME" ]; then
        echo "请设置 Android NDK 路径，例如："
        echo "export ANDROID_NDK_HOME=/path/to/android-ndk"
        exit 1
    fi
    
    echo "Android NDK: $ANDROID_NDK_HOME"
    
    # 获取版本信息
    VERSION="v0.0.1-$(date +%Y%m%d)-$(git rev-parse --short HEAD 2>/dev/null || echo 'dev')"
    echo "版本: $VERSION"
    
    # 创建输出目录
    mkdir -p .build
    
        # 构建 ARM64 版本（真机）
    echo ""
    echo "构建 ARM64 版本..."
    fyne package -os android -appID com.oneclickvirt.goecs -appVersion "$VERSION"
    
    if [ -f *.apk ]; then
        mv *.apk .build/goecs-android-arm64-${VERSION}.apk
        echo "ARM64 APK 构建成功"
    else
        echo "ARM64 APK 构建失败"
        exit 1
    fi
    
    # 构建 x86_64 版本（模拟器）
    echo ""
    echo "构建 x86_64 版本..."
    fyne package -os android/amd64 -appID com.oneclickvirt.goecs -appVersion "$VERSION"
    
    if [ -f *.apk ]; then
        mv *.apk .build/goecs-android-x86_64-${VERSION}.apk
        echo "x86_64 APK 构建成功"
    else
        echo "x86_64 APK 构建失败"
        exit 1
    fi
    
    echo ""
    echo "=========================================="
    echo "  构建完成"
    echo "=========================================="
    ls -lh .build/
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
    "all")
        build_desktop
        echo ""
        check_fyne_cli
        build_android
        ;;
    *)
        echo "用法: $0 [desktop|android|all]"
        echo ""
        echo "  desktop - 构建桌面端应用（默认，用于快速测试）"
        echo "  android - 构建 Android APK"
        echo "  all     - 构建所有平台"
        exit 1
        ;;
esac

echo ""
echo "所有构建任务完成"
