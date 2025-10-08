package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

// ImportTemplate 导入模板文件
func (a *App) ImportTemplate() (*TemplateFile, error) {
	// 打开文件选择对话框
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择模板文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "CLIQ模板文件 (*.yaml, *.yml)",
				Pattern:     "*.yaml;*.yml",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	// 用户取消选择
	if filePath == "" {
		return nil, errors.New("未选择文件")
	}

	// 解析模板文件
	template, err := ParseTemplateFile(filePath)
	if err != nil {
		return nil, err
	}

	a.template = template // 将解析后的模板赋值给 App 结构体

	return template, nil
}

// ParseTemplateFile 解析模板文件
func ParseTemplateFile(filePath string) (*TemplateFile, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析YAML
	var template TemplateFile
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}

	// 验证模板
	if err := validateTemplate(&template); err != nil {
		return nil, err
	}

	return &template, nil
}

// ValidateTemplate 验证模板是否合法
func validateTemplate(template *TemplateFile) error {
	// 验证基本信息
	if template.Name == "" {
		return errors.New("模板名称不能为空")
	}

	if template.CliqTemplateVersion == "" {
		return errors.New("模板版本不能为空")
	}

	// 验证命令
	if len(template.Cmds) == 0 {
		return errors.New("模板必须包含至少一个命令")
	}

	for i, cmd := range template.Cmds {
		if cmd.Name == "" {
			return fmt.Errorf("命令 #%d 名称不能为空", i+1)
		}

		if cmd.Command == "" {
			return fmt.Errorf("命令 #%d 命令模板不能为空", i+1)
		}

		// 验证变量
		for name, variable := range cmd.Variables {
			if variable.Label == "" {
				return fmt.Errorf("命令 #%d 变量 %s 的标签不能为空", i+1, name)
			}

			// 验证变量类型
			switch variable.Type {
			case VarTypeText, VarTypeFileInput, VarTypeFileOutput, VarTypeBoolean, VarTypeNumber, VarTypeSelect:
				// 合法类型
			default:
				return fmt.Errorf("命令 #%d 变量 %s 的类型 %s 不支持", i+1, name, variable.Type)
			}
		}
	}

	return nil
}
