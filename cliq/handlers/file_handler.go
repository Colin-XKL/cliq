package handlers

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"cliq/models"
	"gopkg.in/yaml.v3"
)

// FileHandler handles file-related operations
type FileHandler struct {
	ctx context.Context
}

// NewFileHandler creates a new file handler
func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (fh *FileHandler) Startup(ctx context.Context) {
	fh.ctx = ctx
}

// OpenFileDialog opens a file dialog and returns the selected file path
func (fh *FileHandler) OpenFileDialog() (string, error) {
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

	return runtime.OpenFileDialog(fh.ctx, options)
}

// SaveFileDialog opens a save file dialog and returns the selected file path
func (fh *FileHandler) SaveFileDialog() (string, error) {
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

	return runtime.SaveFileDialog(fh.ctx, options)
}

// SaveYAMLToFile opens a save file dialog and saves the YAML content to the selected file
func (fh *FileHandler) SaveYAMLToFile(yamlContent string) error {
	if yamlContent == "" {
		return fmt.Errorf("YAML内容不能为空")
	}

	// Open save file dialog
	filePath, err := runtime.SaveFileDialog(fh.ctx, runtime.SaveDialogOptions{
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

// ExportTemplateToFile 将模板导出为文件
func (fh *FileHandler) ExportTemplateToFile(template *models.TemplateFile, filePath string) error {
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

func getCommandParts(cmdTemplateStr string, variables map[string]interface{}) ([]string, error) {
	// 替换命令模板中的变量
	commandStr := cmdTemplateStr
	for name, value := range variables {
		// TODO: 根据变量类型进行格式化，例如布尔值转换为 --flag 或空
		commandStr = strings.ReplaceAll(commandStr, fmt.Sprintf("{{%s}}", name), fmt.Sprintf("%v", value))
	}

	parts := strings.Fields(commandStr)
	return parts, nil
}

// ExecuteCommand executes a shell command with the given input and output file paths
func (fh *FileHandler) ExecuteCommand(template *models.TemplateFile, commandID string, variables map[string]interface{}) (string, error) {
	// 根据 commandID 查找对应的命令
	var selectedCommand models.Command
	found := false
	for _, cmd := range template.Cmds {
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
		return "", fmt.Errorf("命令为空")
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

func (fh *FileHandler) GetCommandText(template *models.TemplateFile, commandID string, variables map[string]interface{}) (string, error) {
	var selectedCommand models.Command
	found := false
	for _, cmd := range template.Cmds {
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

// getHashForTemplateName 生成基于模板名称的安全哈希值，防止路径遍历和特殊字符问题
func getHashForTemplateName(templateName string) string {
	hash := md5.Sum([]byte(templateName))
	return fmt.Sprintf("%x", hash)
}

// getFavTemplatesDirPath 获取收藏模板的存储路径
func (fh *FileHandler) getFavTemplatesDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户主目录失败: %w", err)
	}
	favTemplatesDir := filepath.Join(homeDir, ".config", "cliq", "fav_templates")
	return favTemplatesDir, nil
}

// ensureFavTemplatesDirExists 确保收藏模板目录存在，如果不存在则创建
func (fh *FileHandler) ensureFavTemplatesDirExists() (string, error) {
	dirPath, err := fh.getFavTemplatesDirPath()
	if err != nil {
		return "", err
	}

	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return "", fmt.Errorf("创建收藏模板目录失败: %w", err)
		}
	}

	return dirPath, nil
}

// SaveFavTemplate 保存收藏模板文件
func (fh *FileHandler) SaveFavTemplate(template *models.TemplateFile) error {
	if template == nil {
		return fmt.Errorf("模板不能为空")
	}

	// 确保收藏目录存在
	dirPath, err := fh.ensureFavTemplatesDirExists()
	if err != nil {
		return err
	}

	// 使用模板名称的哈希值作为文件名，并添加.cliqfile.yaml后缀
	hashedName := getHashForTemplateName(template.Name)
	fileName := fmt.Sprintf("%s.cliqfile.yaml", hashedName)
	filePath := filepath.Join(dirPath, fileName)

	// 序列化模板为YAML
	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("序列化模板失败: %w", err)
	}

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入收藏模板文件失败: %w", err)
	}

	return nil
}

// ListFavTemplates 列出所有收藏的模板文件
func (fh *FileHandler) ListFavTemplates() ([]*models.TemplateFile, error) {
	dirPath, err := fh.ensureFavTemplatesDirExists()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取收藏模板目录失败: %w", err)
	}

	var templates []*models.TemplateFile
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".cliqfile.yaml") || strings.HasSuffix(file.Name(), ".cliqfile.yml")) {
			filePath := filepath.Join(dirPath, file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("读取文件 %s 失败: %v\n", filePath, err)
				continue
			}

			var template models.TemplateFile
			err = yaml.Unmarshal(data, &template)
			if err != nil {
				fmt.Printf("解析文件 %s 失败: %v\n", filePath, err)
				continue
			}
			templates = append(templates, &template)
		}
	}

	return templates, nil
}

// isValidMD5Hash 验证字符串是否为有效的MD5哈希值格式
func isValidMD5Hash(hash string) bool {
	if len(hash) != 32 {
		return false
	}
	
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	
	return true
}

// getFilePathForTemplate 根据模板名称获取文件路径，同时兼容新旧格式
func (fh *FileHandler) getFilePathForTemplate(templateName string) (string, error) {
	dirPath, err := fh.getFavTemplatesDirPath()
	if err != nil {
		return "", err
	}
	
	// 尝试新格式（使用哈希值）
	hashedName := getHashForTemplateName(templateName)
	newFileName := fmt.Sprintf("%s.cliqfile.yaml", hashedName)
	newFilePath := filepath.Join(dirPath, newFileName)
	
	// 检查新格式的文件是否存在
	if _, err := os.Stat(newFilePath); err == nil {
		return newFilePath, nil
	}
	
	// 如果新格式不存在，检查旧格式的文件是否存在
	oldFileName := fmt.Sprintf("%s.cliqfile.yaml", templateName)
	oldFilePath := filepath.Join(dirPath, oldFileName)
	
	if _, err := os.Stat(oldFilePath); err == nil {
		return oldFilePath, nil
	}
	
	// 都不存在则返回新格式的路径（用于保存等操作）
	return newFilePath, nil
}

// DeleteFavTemplate 从收藏目录删除指定模板文件
func (fh *FileHandler) DeleteFavTemplate(templateName string) error {
	if templateName == "" {
		return fmt.Errorf("模板名称不能为空")
	}

	filePath, err := fh.getFilePathForTemplate(templateName)
	if err != nil {
		return err
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("模板文件不存在: %s", templateName)
	}

	// 删除文件
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("删除收藏模板文件失败: %w", err)
	}

	return nil
}

// GetFavTemplate 读取指定收藏模板文件内容
func (fh *FileHandler) GetFavTemplate(templateName string) (*models.TemplateFile, error) {
	if templateName == "" {
		return nil, fmt.Errorf("模板名称不能为空")
	}

	filePath, err := fh.getFilePathForTemplate(templateName)
	if err != nil {
		return nil, err
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("模板文件不存在: %s", templateName)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取收藏模板文件失败: %w", err)
	}

	var template models.TemplateFile
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return nil, fmt.Errorf("解析收藏模板文件失败: %w", err)
	}

	return &template, nil
}

// UpdateFavTemplate 更新指定收藏模板文件内容
func (fh *FileHandler) UpdateFavTemplate(templateName string, updatedTemplate *models.TemplateFile) error {
	if templateName == "" {
		return fmt.Errorf("模板名称不能为空")
	}
	if updatedTemplate == nil {
		return fmt.Errorf("更新模板不能为空")
	}

	filePath, err := fh.getFilePathForTemplate(templateName)
	if err != nil {
		return err
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("模板文件不存在: %s", templateName)
	}

	// 序列化更新后的模板为YAML
	data, err := yaml.Marshal(updatedTemplate)
	if err != nil {
		return fmt.Errorf("序列化更新模板失败: %w", err)
	}

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入收藏模板文件失败: %w", err)
	}

	return nil
}