package main

import (
	"context"
	"fmt"
	"os/exec"
	goruntime "runtime"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// OpenFileDialog opens a file dialog and returns the selected file path
func (a *App) OpenFileDialog() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "选择输入文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PNG图片 (*.png)",
				Pattern:     "*.png",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	}

	return runtime.OpenFileDialog(a.ctx, options)
}

// SaveFileDialog opens a save file dialog and returns the selected file path
func (a *App) SaveFileDialog() (string, error) {
	options := runtime.SaveDialogOptions{
		Title: "保存输出文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PNG图片 (*.png)",
				Pattern:     "*.png",
			},
		},
		DefaultFilename: "output.png",
	}

	return runtime.SaveFileDialog(a.ctx, options)
}

// ExecuteCommand executes a shell command with the given input and output file paths
func (a *App) ExecuteCommand(inputFilePath string, outputFilePath string) (string, error) {
	var cmd *exec.Cmd

	// 构建命令
	commandStr := fmt.Sprintf("pngquant %s --output %s", inputFilePath, outputFilePath)

	// 根据操作系统选择不同的命令执行方式
	if goruntime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", commandStr)
	} else {
		cmd = exec.Command("sh", "-c", commandStr)
	}

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %v, 输出: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
