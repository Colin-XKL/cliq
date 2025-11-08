package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"

	"cliq/models"
)

// ImportTemplate 导入模板文件
func (a *App) ImportTemplate() (*models.TemplateFile, error) {
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

// ImportTemplateFromURL 从URL导入模板文件
func (a *App) ImportTemplateFromURL(url string) (*models.TemplateFile, error) {
	// 从URL下载内容
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("下载模板文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取模板内容失败: %w", err)
	}

	// 解析YAML
	var template models.TemplateFile
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}

	// 验证模板
	if err := validateTemplate(&template); err != nil {
		return nil, err
	}

	a.template = &template // 将解析后的模板赋值给 App 结构体

	return &template, nil
}

// ParseTemplateFile 解析模板文件
func ParseTemplateFile(filePath string) (*models.TemplateFile, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析YAML
	var template models.TemplateFile
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}

	// 验证模板（包括变量名唯一性）
	if err := validateTemplate(&template); err != nil {
		return nil, err
	}

	return &template, nil
}

// validateVariableNamesInTemplate 验证模板中变量名的唯一性
func validateVariableNamesInTemplate(template *models.TemplateFile) error {
	for i, cmd := range template.Cmds {
		seen := make(map[string]bool)
		for _, varDef := range cmd.Variables {
			if varDef.Name == "" {
				return fmt.Errorf("命令 #%d 变量名称不能为空", i+1)
			}
			if seen[varDef.Name] {
				return fmt.Errorf("命令 #%d 中存在重复变量名: %s", i+1, varDef.Name)
			}
			seen[varDef.Name] = true
		}
	}
	return nil
}

// ValidateTemplate 验证模板是否合法
func validateTemplate(template *models.TemplateFile) error {
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
		for _, varDef := range cmd.Variables {
			if varDef.Name == "" {
				return fmt.Errorf("命令 #%d 变量名称不能为空", i+1)
			}
			if varDef.Label == "" {
				return fmt.Errorf("命令 #%d 变量 %s 的标签不能为空", i+1, varDef.Name)
			}

			// 验证变量类型
			switch varDef.Type {
			case models.VarTypeText, models.VarTypeFileInput, models.VarTypeFileOutput, models.VarTypeBoolean, models.VarTypeNumber, models.VarTypeSelect:
				// 合法类型
			default:
				return fmt.Errorf("命令 #%d 变量 %s 的类型 %s 不支持", i+1, varDef.Name, varDef.Type)
			}
		}
	}

	// 验证变量名唯一性
	if err := validateVariableNamesInTemplate(template); err != nil {
		return err
	}

	return nil
}
