# cliQ

[English](./README.md) | [中文](./README_zh-CN.md)

**cliQ** 发音：/klɪk/，类似“click”
cliQ 是一个轻量级工具，旨在将复杂的 CLI 命令转化为直观、易用的图形用户界面（GUI）。

用户只需定义一条带有变量占位符的命令模板，即可自动生成对应的 GUI 界面，通过点击、选择等方式完成参数输入，无需记忆命令语法，极大提升使用体验和操作效率。

## 主要功能
- 命令模板定义：用户可输入 CLI 命令，生成对应的动态表单界面。
- 模板导入/导出：支持将模板(`.cliqfile.yaml`)导出为文件，或从他人导入，便于团队共享。
- 多种输入组件支持：支持文件选择器、数字输入框、下拉选择框、复选框等，适配不同参数类型。
- 跨平台支持：支持 Windows、macOS 和 Linux 平台
- 模板市场：用户可上传/下载常用工具模板（如 ImageMagick、ffmpeg、pngquant 等），构建共享生态。

## 如何使用

从Release界面下载对应平台的安装包，双击运行即可.

**Windows 平台**:
exe文件直接打开运行即可. Windows 10及以上版本均支持.
程序依赖Webview, 如果本机没有安装, 第一次打开时会弹出提示,按照指引安装即可.

**macOS 平台**:
macOS 10.13及以上版本均支持.
由于产物没有签名, 直接打开会弹出提示“文件已损坏”, 请按照以下步骤解决:
1. 打开终端, 手动认证cliQ应用:
```bash
sudo xattr -dr com.apple.quarantine /Users/colin/Downloads/cliq.app
# 注: 最后面的路径请根据实际情况修改. 你也可以直接将文件拖入到终端, 会自动补齐路径.
```
2. 打开cliQ应用, 如果能正常打开, 则说明认证成功.
3. (可选)将cliQ应用拖动到“应用程序”文件夹, 方便后续使用.

**Linux平台**:
cliQ应用本身支持Linux平台运行. 当前产物尚未在Linux平台充分测试, 欢迎通过issue反馈.

## 技术栈
- 使用wails框架构建, 跨平台支持
- 页面使用Vue3, 使用TypeScript
- 使用PrimeVue V4组件库、Vite构建前端项目, 部分样式使用TailwindCSS