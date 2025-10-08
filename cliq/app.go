package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
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

// ParseCommandToTemplate 将命令字符串解析为模板
func (a *App) ParseCommandToTemplate(commandStr string) (*TemplateFile, error) {
	if commandStr == "" {
		return nil, fmt.Errorf("命令字符串不能为空")
	}
	
	// 从命令字符串中提取变量
	variables := extractVariablesFromCommand(commandStr)
	
	// 生成模板
	templateFile := &TemplateFile{
		Name:                "Generated Template",
		Description:         "Automatically generated template from command",
		Version:             "1.0",
		Author:              "cliQ",
		CliqTemplateVersion: "1.0",
		Cmds: []Command{
			{
				ID:          "generated_cmd_1",
				Name:        "Generated Command",
				Description: "Automatically generated command",
				Command:     commandStr,
				Variables:   map[string]Variable{},
			},
		},
	}
	
	// 为每个提取的变量创建适当的参数配置
	for _, varName := range variables {
		varType := determineVariableType(varName)
		variable := Variable{
			Type:        varType,
			Label:       getLabelFromVariableName(varName),
			Description: fmt.Sprintf("The %s parameter", varName),
			Required:    true,
		}
		
		// 根据变量类型设置特定选项
		if varType == VarTypeFileInput || varType == VarTypeFileOutput {
			variable.Options = map[string]interface{}{
				"file_types": []string{".*"}, // 默认支持所有文件类型
			}
		} else if varType == VarTypeNumber {
			variable.Options = map[string]interface{}{
				"default": 1,
				"min":     0,
				"max":     100,
			}
		} else if varType == VarTypeBoolean {
			variable.Options = map[string]interface{}{
				"default": false,
			}
		}
		
		templateFile.Cmds[0].Variables[varName] = variable
	}
	
	return templateFile, nil
}

// GenerateYAMLFromTemplate 将模板对象转换为YAML字符串
func (a *App) GenerateYAMLFromTemplate(template *TemplateFile) (string, error) {
	if template == nil {
		return "", fmt.Errorf("模板不能为空")
	}

	// 序列化模板为YAML
	data, err := yaml.Marshal(template)
	if err != nil {
		return "", fmt.Errorf("序列化模板失败: %w", err)
	}

	return string(data), nil
}

// ValidateYAMLTemplate 验证YAML模板格式
func (a *App) ValidateYAMLTemplate(yamlStr string) error {
	if yamlStr == "" {
		return fmt.Errorf("YAML字符串不能为空")
	}

	// 反序列化YAML到TemplateFile结构
	var template TemplateFile
	err := yaml.Unmarshal([]byte(yamlStr), &template)
	if err != nil {
		return fmt.Errorf("YAML格式错误: %w", err)
	}

	// 验证模板结构
	return validateTemplate(&template)
}

// extractVariablesFromCommand 从命令字符串中提取变量名
func extractVariablesFromCommand(commandStr string) []string {
	var variables []string
	parts := strings.Fields(commandStr)
	
	for _, part := range parts {
		// 查找 {{variable_name}} 格式的变量
		for start := 0; start < len(part); {
			pos := strings.Index(part[start:], "{{")
			if pos == -1 {
				break
			}
			start += pos
			
			endPos := strings.Index(part[start:], "}}")
			if endPos == -1 {
				break
			}
			endPos += start
			
			variable := strings.TrimSpace(part[start+2:endPos])
			// 检查是否已存在
			exists := false
			for _, v := range variables {
				if v == variable {
					exists = true
					break
				}
			}
			if !exists {
				variables = append(variables, variable)
			}
			
			start = endPos + 2
		}
	}
	
	return variables
}

// determineVariableType 根据变量名确定变量类型
func determineVariableType(varName string) string {
	// 根据变量名后缀判断类型
	if strings.HasSuffix(varName, "_file") || strings.HasSuffix(varName, "_path") || 
	   strings.Contains(varName, "file") || strings.Contains(varName, "path") {
		// 判断是输入文件还是输出文件
		if strings.Contains(varName, "input") || strings.Contains(varName, "src") {
			return VarTypeFileInput
		} else if strings.Contains(varName, "output") || strings.Contains(varName, "dest") {
			return VarTypeFileOutput
		} else {
			// 如果只是file或path，缺省为输入文件
			return VarTypeFileInput
		}
	} else if strings.Contains(varName, "number") || strings.Contains(varName, "size") || 
	          strings.Contains(varName, "width") || strings.Contains(varName, "height") {
		return VarTypeNumber
	} else if strings.Contains(varName, "enable") || strings.Contains(varName, "use") || 
	          strings.Contains(varName, "flag") || strings.Contains(varName, "show") {
		return VarTypeBoolean
	} else {
		return VarTypeText
	}
}

// getLabelFromVariableName 根据变量名生成标签
func getLabelFromVariableName(varName string) string {
	// 将变量名转换为更友好的标签格式
	label := strings.ReplaceAll(varName, "_", " ")
	label = strings.ReplaceAll(label, "-", " ")
	// 首字母大写
	if len(label) > 0 {
		label = strings.ToUpper(string(label[0])) + label[1:]
	}
	return label
}

// ExportTemplateToFile 将模板导出为文件
func (a *App) ExportTemplateToFile(template *TemplateFile, filePath string) error {
	if template == nil {
		return fmt.Errorf("模板不能为空")
	}
	
	if filePath == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	// 序列化模板为YAML
	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("序列化模板失败: %w", err)
	}

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// SaveYAMLToFile opens a save file dialog and saves the YAML content to the selected file
func (a *App) SaveYAMLToFile(yamlContent string) error {
	if yamlContent == "" {
		return fmt.Errorf("YAML内容不能为空")
	}

	// Open save file dialog
	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "保存模板文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "CLIQ模板文件 (*.yaml, *.yml)",
				Pattern:     "*.yaml;*.yml",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
		DefaultFilename: "template.cliqfile.yaml",
	})

	if err != nil {
		return fmt.Errorf("打开保存对话框失败: %w", err)
	}

	// 用户取消选择
	if filePath == "" {
		return fmt.Errorf("未选择保存路径")
	}

	// Write content to file
	err = os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
