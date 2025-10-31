# Android APK 构建产物

此目录存放 GitHub Actions 自动构建的 APK 文件。

## 文件命名规则

`goecs-android-{arch}-{version}.apk`

- arch: arm64 或 x86_64
- version: 版本号格式为 v0.0.1-YYYYMMDD-{git-hash}

## 架构说明

- **arm64**: 适用于真实 Android 设备（推荐）
- **x86_64**: 适用于 Android 模拟器

## 使用方法

下载对应架构的 APK 文件，传输到 Android 设备上安装即可。

