package services

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"cliq/models"
)

// TemplateService provides template-related business logic
type TemplateService struct{}

// NewTemplateService creates a new template service
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// ParseCommandToTemplate 将命令字符串解析为模板
func (ts *TemplateService) ParseCommandToTemplate(commandStr string) (*models.TemplateFile, error) {
	if commandStr == "" {
		return nil, fmt.Errorf("命令字符串不能为空")
	}

	// 从命令字符串中提取变量
	variables := extractVariablesFromCommand(commandStr)

	// 生成模板
	templateFile := &models.TemplateFile{
		Name:                "Generated Template",
		Description:         "Automatically generated template from command",
		Version:             "1.0",
		Author:              "cliQ",
		CliqTemplateVersion: "1.0",
		Cmds: []models.Command{
			{
				ID:          "generated_cmd_1",
				Name:        "Generated Command",
				Description: "Automatically generated command",
				Command:     commandStr,
				Variables:   map[string]models.Variable{},
			},
		},
	}

	// 为每个提取的变量创建适当的参数配置
	for _, varName := range variables {
		varType := determineVariableType(varName)
		variable := models.Variable{
			Type:        varType,
			Label:       getLabelFromVariableName(varName),
			Description: fmt.Sprintf("The %s parameter", varName),
			Required:    true,
		}

		// 根据变量类型设置特定选项
		switch varType {
		case models.VarTypeFileInput, models.VarTypeFileOutput:
			variable.Options = map[string]interface{}{
				"file_types": []string{".*"}, // 默认支持所有文件类型
			}
		case models.VarTypeNumber:
			variable.Options = map[string]interface{}{
				"default": 1,
				"min":     0,
				"max":     100,
			}
		case models.VarTypeBoolean:
			variable.Options = map[string]interface{}{
				"default": false,
			}
		}

		templateFile.Cmds[0].Variables[varName] = variable
	}

	return templateFile, nil
}

// GenerateYAMLFromTemplate 将模板对象转换为YAML字符串
func (ts *TemplateService) GenerateYAMLFromTemplate(template *models.TemplateFile) (string, error) {
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
func (ts *TemplateService) ValidateYAMLTemplate(yamlStr string) error {
	if yamlStr == "" {
		return fmt.Errorf("YAML字符串不能为空")
	}

	// 反序列化YAML到TemplateFile结构
	var template models.TemplateFile
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

			variable := strings.TrimSpace(part[start+2 : endPos])
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
			return models.VarTypeFileInput
		} else if strings.Contains(varName, "output") || strings.Contains(varName, "dest") {
			return models.VarTypeFileOutput
		} else {
			// 如果只是file或path，缺省为输入文件
			return models.VarTypeFileInput
		}
	} else if strings.Contains(varName, "number") || strings.Contains(varName, "size") ||
		strings.Contains(varName, "width") || strings.Contains(varName, "height") {
		return models.VarTypeNumber
	} else if strings.Contains(varName, "enable") || strings.Contains(varName, "use") ||
		strings.Contains(varName, "flag") || strings.Contains(varName, "show") {
		return models.VarTypeBoolean
	} else {
		return models.VarTypeText
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

// validateTemplate 验证模板是否合法
func validateTemplate(template *models.TemplateFile) error {
	// 验证基本信息
	if template.Name == "" {
		return fmt.Errorf("模板名称不能为空")
	}

	if template.CliqTemplateVersion == "" {
		return fmt.Errorf("模板版本不能为空")
	}

	// 验证命令
	if len(template.Cmds) == 0 {
		return fmt.Errorf("模板必须包含至少一个命令")
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
			case models.VarTypeText, models.VarTypeFileInput, models.VarTypeFileOutput, models.VarTypeBoolean, models.VarTypeNumber, models.VarTypeSelect:
				// 合法类型
			default:
				return fmt.Errorf("命令 #%d 变量 %s 的类型 %s 不支持", i+1, name, variable.Type)
			}
		}
	}

	return nil
}
