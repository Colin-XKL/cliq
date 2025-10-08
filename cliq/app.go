package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	template *TemplateFile // 添加 template 字段
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// TemplateFile 表示一个完整的模板文件
type TemplateFile struct {
	// 模板元信息
	Name                string `yaml:"name" json:"name"`
	Description         string `yaml:"description" json:"description"`
	Version             string `yaml:"version" json:"version"`
	Author              string `yaml:"author" json:"author"`
	CliqTemplateVersion string `yaml:"cliq_template_version" json:"cliq_template_version"`

	// 命令列表
	Cmds []Command `yaml:"cmds" json:"cmds"`
}

// Command 表示一个命令模板
type Command struct {
	ID          string              `yaml:"id" json:"id"` // 添加 ID 字段
	Name        string              `yaml:"name" json:"name"`
	Description string              `yaml:"description" json:"description"`
	Command     string              `yaml:"command" json:"command"`
	Variables   map[string]Variable `yaml:"variables" json:"variables"`
}

// Variable 表示命令中的一个变量
type Variable struct {
	Type        string                 `yaml:"type" json:"type"`
	ArgName     string                 `yaml:"arg_name,omitempty" json:"arg_name,omitempty"`
	Label       string                 `yaml:"label" json:"label"`
	Description string                 `yaml:"description" json:"description"`
	Required    bool                   `yaml:"required" json:"required"`
	Options     map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

// 变量类型常量
const (
	VarTypeText       = "text"
	VarTypeFileInput  = "file_input"
	VarTypeFileOutput = "file_output"
	VarTypeBoolean    = "boolean"
	VarTypeNumber     = "number"
	VarTypeSelect     = "select"
)

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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
func (a *App) ExecuteCommand(commandID string, variables map[string]interface{}) (string, error) {
	// 根据 commandID 查找对应的命令

	var selectedCommand Command
	found := false
	for _, cmd := range a.template.Cmds {
		if cmd.ID == commandID {
			selectedCommand = cmd
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("未找到命令: %s", commandID)
	}

	// 替换命令模板中的变量
	cmdTemplateStr := selectedCommand.Command
	parts, err := getCommandParts(cmdTemplateStr, variables)
	if err != nil {
		return "", fmt.Errorf("获取命令文本失败: %w", err)
	}
	if len(parts) == 0 {
		return "", errors.New("命令为空")
	}

	name := parts[0]
	args := parts[1:]

	cmd := exec.Command(name, args...)

	// 执行命令并获取输出
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %w\n%s", err, string(out))
	}

	return strings.TrimSpace(string(out)), nil
}

func (a *App) GetCommandText(commandID string, variables map[string]interface{}) (string, error) {
	var selectedCommand Command
	found := false
	for _, cmd := range a.template.Cmds {
		if cmd.ID == commandID {
			selectedCommand = cmd
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("未找到命令: %s", commandID)
	}

	// 替换命令模板中的变量
	cmdTemplateStr := selectedCommand.Command
	parts, err := getCommandParts(cmdTemplateStr, variables)
	if err != nil {
		return "", fmt.Errorf("获取命令文本失败: %w", err)
	}
	return strings.Join(parts, " "), nil
}

func getCommandParts(cmdTemplateStr string, variables map[string]interface{}) ([]string, error) {
	// 替换命令模板中的变量
	commandStr := cmdTemplateStr
	for name, value := range variables {
		fmt.Println(name, value)
		// TODO: 根据变量类型进行格式化，例如布尔值转换为 --flag 或空
		commandStr = strings.ReplaceAll(commandStr, fmt.Sprintf("{{%s}}", name), fmt.Sprintf("%v", value))
	}

	parts := strings.Fields(commandStr)
	return parts, nil
}
