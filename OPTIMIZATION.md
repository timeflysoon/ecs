# Android APK 体积优化说明

## 问题分析

当前 Android APK 体积过大的主要原因：

### Vendor 目录大小分析
```
总计: 213MB
├── github.com (179MB)
│   └── oneclickvirt (146MB)
│       ├── fio (101MB) - 包含多平台二进制文件
│       ├── dd (28MB) - 包含多平台二进制文件
│       ├── stream (8.1MB)
│       └── mbw (7.5MB)
├── golang.org (23MB)
└── fyne.io (10MB)
```

### fio 包内容
包含了所有平台的预编译二进制文件：
- fio-linux-386 (5.9MB)
- fio-linux-arm64 (5.9MB)
- fio-linux-amd64 (6.2MB)
- fio-linux-mips64 (6.6MB)
- fio-windows-386.exe (5.7MB)
- fio-windows-amd64.exe (6.2MB)
- fio-darwin-amd64 (1.5MB)
- fio-freebsd-amd64 (11MB)
- 等其他平台...

## 已实施的优化措施

### 1. 添加 ldflags 编译优化
```bash
--ldflags "-s -w -checklinkname=0"
```
- `-s`: 去除符号表
- `-w`: 去除 DWARF 调试信息
- `-checklinkname=0`: 跳过链接名称检查

### 2. 启用 release 模式
```bash
--release
```
这会启用额外的优化和压缩。

### 3. 清理 vendor 中不需要的二进制文件
在构建前删除：
- Darwin (macOS) 平台文件
- Windows 平台文件
- FreeBSD 平台文件
- 32位 Linux (386) 文件
- MIPS/PPC64/S390X 架构文件
- 所有 .exe 文件

只保留 Android 需要的 arm64 和 amd64 二进制文件。

### 4. 创建 .fyneignore 文件
排除不必要的文件：
- 测试文件
- 文档文件
- 示例代码
- 不需要的平台二进制

## 预期效果

根据优化措施，预计可以减少：

1. **ldflags 优化**: 减少约 10-20%
2. **清理 vendor 二进制**: 减少约 80MB (删除不需要的平台文件)
3. **release 模式**: 额外减少 5-10%

**预计最终 APK 大小**: 
- ARM64: 从 223MB → 约 60-80MB
- x86_64: 从 59.2MB → 约 30-40MB

## 进一步优化建议

### 1. 按需加载二进制文件
修改代码逻辑，在运行时从远程下载需要的二进制文件，而不是打包到 APK 中。

### 2. 拆分 APK
使用 Android App Bundle (AAB) 格式：
```bash
fyne package --os android --app-bundle
```
这样可以为不同架构生成不同的 APK。

### 3. 代码层面优化
考虑：
- 移除不需要的依赖
- 使用更轻量的替代库
- 延迟加载重型模块

### 4. 资源优化
- 压缩图片资源
- 移除未使用的资源
- 使用 WebP 格式图片

## 如何验证

下次构建后，查看构建日志中的文件大小输出：
```bash
ls -lh .build/goecs-gui-android-arm64-*.apk
```

## 相关文件

- `.github/workflows/build.yml` - 构建配置（已优化）
- `.fyneignore` - 打包时排除的文件列表（新增）
- `go.mod` - 依赖管理
